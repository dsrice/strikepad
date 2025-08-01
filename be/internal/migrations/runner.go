package migrations

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"ariga.io/atlas-go-sdk/atlasexec"
)

type MigrationRunner struct {
	client *atlasexec.Client
	env    string
}

// NewMigrationRunner creates a new migration runner instance
func NewMigrationRunner(env string) (*MigrationRunner, error) {
	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Set Atlas binary path
	atlasPath := filepath.Join(workingDir, "atlas.exe")

	// Check if Atlas binary exists
	if _, err := os.Stat(atlasPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("atlas binary not found at %s", atlasPath)
	}

	// Create Atlas client
	client, err := atlasexec.NewClient(workingDir, atlasPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create atlas client: %w", err)
	}

	return &MigrationRunner{
		client: client,
		env:    env,
	}, nil
}

// RunMigrations executes pending migrations
func (mr *MigrationRunner) RunMigrations(ctx context.Context) error {
	slog.Info("Running database migrations", "environment", mr.env)

	// Apply pending migrations
	result, err := mr.client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		Env: mr.env,
	})
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if len(result.Applied) > 0 {
		slog.Info("Successfully applied migrations", "count", len(result.Applied))
		for _, applied := range result.Applied {
			slog.Debug("Applied migration", "name", applied.Name)
		}
	} else {
		slog.Info("No pending migrations found")
	}

	return nil
}

// ValidateMigrations validates the migration files
func (mr *MigrationRunner) ValidateMigrations(ctx context.Context) error {
	slog.Info("Validating migrations", "environment", mr.env)

	// Get migration status to validate
	_, err := mr.client.MigrateStatus(ctx, &atlasexec.MigrateStatusParams{
		Env: mr.env,
	})
	if err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	slog.Info("Migration validation passed")
	return nil
}

// GetMigrationStatus returns the current migration status
func (mr *MigrationRunner) GetMigrationStatus(ctx context.Context) (*atlasexec.MigrateStatus, error) {
	return mr.client.MigrateStatus(ctx, &atlasexec.MigrateStatusParams{
		Env: mr.env,
	})
}
