package container

import (
	"strikepad-manage-tool/config"
	"strikepad-manage-tool/handlers"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/repository/ri"
	"strikepad-manage-tool/templates"

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

func NewUserRepository(db *gorm.DB) ri.UserRepositoryInterface {
	return repository.NewUserRepository(db)
}

func NewAdminRepository(db *gorm.DB) ri.AdminRepositoryInterface {
	return repository.NewAdminRepository(db)
}

func NewMakerRepository(db *gorm.DB) ri.MakerRepositoryInterface {
	return repository.NewMakerRepository(db)
}

func NewCoreRepository(db *gorm.DB) ri.CoreRepositoryInterface {
	return repository.NewCoreRepository(db)
}

func NewCoverRepository(db *gorm.DB) ri.CoverRepositoryInterface {
	return repository.NewCoverRepository(db)
}

// ハンドラー層のコンストラクター関数

func NewAuthHandler(adminRepo ri.AdminRepositoryInterface) handlers.AuthHandlerInterface {
	return handlers.NewAuthHandler(adminRepo)
}

func NewAdminHandler(userRepo ri.UserRepositoryInterface) handlers.AdminHandlerInterface {
	return handlers.NewAdminHandler(userRepo)
}

func NewMakerHandler(makerRepo ri.MakerRepositoryInterface, s3Client *config.S3Client) handlers.MakerHandlerInterface {
	return handlers.NewMakerHandler(makerRepo, s3Client)
}

func NewCoreHandler(coreRepo ri.CoreRepositoryInterface, makerRepo ri.MakerRepositoryInterface) handlers.CoreHandlerInterface {
	return handlers.NewCoreHandler(coreRepo, makerRepo)
}

func NewCoverHandler(coverRepo ri.CoverRepositoryInterface, makerRepo ri.MakerRepositoryInterface) handlers.CoverHandlerInterface {
	return handlers.NewCoverHandler(coverRepo, makerRepo)
}

// HandlersContainer は全てのハンドラーを格納する構造体
type HandlersContainer struct {
	AuthHandler  handlers.AuthHandlerInterface
	AdminHandler handlers.AdminHandlerInterface
	MakerHandler handlers.MakerHandlerInterface
	CoreHandler  handlers.CoreHandlerInterface
	CoverHandler handlers.CoverHandlerInterface
}

// NewHandlersContainer は全てのハンドラーを含むコンテナーを作成
func NewHandlersContainer(
	authHandler handlers.AuthHandlerInterface,
	adminHandler handlers.AdminHandlerInterface,
	makerHandler handlers.MakerHandlerInterface,
	coreHandler handlers.CoreHandlerInterface,
	coverHandler handlers.CoverHandlerInterface,
) *HandlersContainer {
	return &HandlersContainer{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		MakerHandler: makerHandler,
		CoreHandler:  coreHandler,
		CoverHandler: coverHandler,
	}
}