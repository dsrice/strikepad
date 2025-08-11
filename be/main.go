package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"strikepad-backend/internal/container"
	"strikepad-backend/internal/handler"
	"strikepad-backend/internal/migrations"
	authMiddleware "strikepad-backend/internal/middleware"
	"strikepad-backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Initialize structured logger
	initLogger()

	// Run database migrations on startup
	if err := runMigrations(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	c := container.BuildContainer()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from StrikePad Backend!")
	})

	err := c.Invoke(
		func(
			healthHandler handler.HealthHandlerInterface,
			apiHandler *handler.APIHandler,
			authHandler handler.AuthHandlerInterface,
			sessionService service.SessionServiceInterface,
		) {
			e.GET("/health", healthHandler.Check)
			e.GET("/api/test", apiHandler.Test)

			// Public auth endpoints (no JWT required)
			e.POST("/api/auth/signup", authHandler.Signup)
			e.POST("/api/auth/login", authHandler.Login)
			e.POST("/api/auth/google/signup", authHandler.GoogleSignup)
			e.POST("/api/auth/google/login", authHandler.GoogleLogin)

			// Protected auth endpoints (JWT required)
			protected := e.Group("/api/auth", authMiddleware.JWTMiddleware(sessionService))
			protected.POST("/logout", authHandler.Logout)
		})

	if err != nil {
		slog.Error("Failed to invoke handlers", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting server", "port", 8080)
	if err := e.Start(":8080"); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// initLogger initializes the structured logger with file output and rotation
func initLogger() {
	// Get log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	var level slog.Level

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0750); err != nil {
		slog.Error("Failed to create logs directory", "error", err)
		os.Exit(1)
	}

	// Setup lumberjack for log rotation (hourly rotation)
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logsDir, "app.log"),
		MaxSize:    100, // MB
		MaxBackups: 24,  // Keep 24 hours of logs
		MaxAge:     7,   // Keep logs for 7 days
		Compress:   true,
	}

	// Create combined writer for both file and stdout
	var writer io.Writer
	env := os.Getenv("APP_ENV")

	if env == "production" {
		// Production: only write to file
		writer = logFile
	} else {
		// Development: write to both file and stdout
		writer = io.MultiWriter(logFile, os.Stdout)
	}

	// Setup handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // Add source file and line number
	}

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Setup hourly log rotation using a goroutine
	setupHourlyRotation(logFile)
}

// runMigrations executes database migrations on application startup
func runMigrations() error {
	// Get environment (default to "dev")
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Skip migrations in test environment to avoid conflicts
	if env == "test" {
		slog.Info("Skipping migrations for test environment")
		return nil
	}

	slog.Info("Initializing migration runner", "environment", env)

	// Create migration runner
	runner, err := migrations.NewMigrationRunner(env)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Run migrations
	if err := runner.RunMigrations(ctx); err != nil {
		return err
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

// setupHourlyRotation sets up hourly log rotation
func setupHourlyRotation(logFile *lumberjack.Logger) {
	go func() {
		// Calculate time until next hour
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		timeUntilNextHour := nextHour.Sub(now)

		// Wait until the next hour
		time.Sleep(timeUntilNextHour)

		// Create ticker for hourly rotation
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			// Force rotation
			if err := logFile.Rotate(); err != nil {
				slog.Error("Failed to rotate log file", "error", err)
			} else {
				slog.Info("Log file rotated successfully")
			}
		}
	}()
}
