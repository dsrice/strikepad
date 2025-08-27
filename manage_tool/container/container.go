package container

import (
	"strikepad-manage-tool/config"
	"strikepad-manage-tool/handlers"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/templates"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// Container はDI用のコンテナー
type Container struct {
	*dig.Container
}

// New は新しいDIコンテナーを作成
func New() *Container {
	container := dig.New()
	return &Container{Container: container}
}

// BuildContainer は依存関係を設定したコンテナーを構築
func BuildContainer() (*Container, error) {
	container := New()

	// データベース関連
	if err := container.Provide(config.NewDatabaseConfig); err != nil {
		return nil, err
	}

	if err := container.Provide(config.ConnectDatabase); err != nil {
		return nil, err
	}

	// S3クライアント
	if err := container.Provide(config.NewS3Client); err != nil {
		return nil, err
	}

	// テンプレートレンダラー
	if err := container.Provide(templates.NewTemplateRenderer); err != nil {
		return nil, err
	}

	// リポジトリ層
	if err := container.Provide(NewUserRepository); err != nil {
		return nil, err
	}

	if err := container.Provide(NewAdminRepository); err != nil {
		return nil, err
	}

	if err := container.Provide(NewMakerRepository); err != nil {
		return nil, err
	}

	if err := container.Provide(NewCoreRepository); err != nil {
		return nil, err
	}

	if err := container.Provide(NewCoverRepository); err != nil {
		return nil, err
	}

	// ハンドラー層
	if err := container.Provide(NewAuthHandler); err != nil {
		return nil, err
	}

	if err := container.Provide(NewAdminHandler); err != nil {
		return nil, err
	}

	if err := container.Provide(NewMakerHandler); err != nil {
		return nil, err
	}

	if err := container.Provide(NewCoreHandler); err != nil {
		return nil, err
	}

	if err := container.Provide(NewCoverHandler); err != nil {
		return nil, err
	}

	// ハンドラーコンテナー
	if err := container.Provide(NewHandlersContainer); err != nil {
		return nil, err
	}

	return container, nil
}

// リポジトリ層のコンストラクター関数

func NewUserRepository(db *gorm.DB) *repository.UserRepository {
	return repository.NewUserRepository(db)
}

func NewAdminRepository(db *gorm.DB) *repository.AdminRepository {
	return repository.NewAdminRepository(db)
}

func NewMakerRepository(db *gorm.DB) *repository.MakerRepository {
	return repository.NewMakerRepository(db)
}

func NewCoreRepository(db *gorm.DB) *repository.CoreRepository {
	return repository.NewCoreRepository(db)
}

func NewCoverRepository(db *gorm.DB) *repository.CoverRepository {
	return repository.NewCoverRepository(db)
}

// ハンドラー層のコンストラクター関数

func NewAuthHandler(adminRepo *repository.AdminRepository) *handlers.AuthHandler {
	return handlers.NewAuthHandler(adminRepo)
}

func NewAdminHandler(userRepo *repository.UserRepository) *handlers.AdminHandler {
	return handlers.NewAdminHandler(userRepo)
}

func NewMakerHandler(makerRepo *repository.MakerRepository, s3Client *s3.Client) *handlers.MakerHandler {
	return handlers.NewMakerHandler(makerRepo, s3Client)
}

func NewCoreHandler(coreRepo *repository.CoreRepository, makerRepo *repository.MakerRepository) *handlers.CoreHandler {
	return handlers.NewCoreHandler(coreRepo, makerRepo)
}

func NewCoverHandler(coverRepo *repository.CoverRepository, makerRepo *repository.MakerRepository) *handlers.CoverHandler {
	return handlers.NewCoverHandler(coverRepo, makerRepo)
}

// HandlersContainer は全てのハンドラーを格納する構造体
type HandlersContainer struct {
	AuthHandler  *handlers.AuthHandler
	AdminHandler *handlers.AdminHandler
	MakerHandler *handlers.MakerHandler
	CoreHandler  *handlers.CoreHandler
	CoverHandler *handlers.CoverHandler
}

// NewHandlersContainer は全てのハンドラーを含むコンテナーを作成
func NewHandlersContainer(
	authHandler *handlers.AuthHandler,
	adminHandler *handlers.AdminHandler,
	makerHandler *handlers.MakerHandler,
	coreHandler *handlers.CoreHandler,
	coverHandler *handlers.CoverHandler,
) *HandlersContainer {
	return &HandlersContainer{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		MakerHandler: makerHandler,
		CoreHandler:  coreHandler,
		CoverHandler: coverHandler,
	}
}