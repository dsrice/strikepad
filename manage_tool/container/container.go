package container

import (
	"strikepad-manage-tool/config"
	"strikepad-manage-tool/handlers"
	"strikepad-manage-tool/handlers/hi"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/templates"

	"go.uber.org/dig"
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
	if err := container.Provide(repository.NewUserRepositoryInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(repository.NewAdminRepositoryInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(repository.NewMakerRepositoryInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(repository.NewCoreRepositoryInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(repository.NewCoverRepositoryInterface); err != nil {
		return nil, err
	}

	// ハンドラー層
	if err := container.Provide(handlers.NewAuthHandlerInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(handlers.NewAdminHandlerInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(handlers.NewMakerHandlerInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(handlers.NewCoreHandlerInterface); err != nil {
		return nil, err
	}

	if err := container.Provide(handlers.NewCoverHandlerInterface); err != nil {
		return nil, err
	}

	// ハンドラーコンテナー
	if err := container.Provide(NewHandlersContainer); err != nil {
		return nil, err
	}

	return container, nil
}


// HandlersContainer は全てのハンドラーを格納する構造体
type HandlersContainer struct {
	AuthHandler  hi.AuthHandlerInterface
	AdminHandler hi.AdminHandlerInterface
	MakerHandler hi.MakerHandlerInterface
	CoreHandler  hi.CoreHandlerInterface
	CoverHandler hi.CoverHandlerInterface
}

// NewHandlersContainer は全てのハンドラーを含むコンテナーを作成
func NewHandlersContainer(
	authHandler hi.AuthHandlerInterface,
	adminHandler hi.AdminHandlerInterface,
	makerHandler hi.MakerHandlerInterface,
	coreHandler hi.CoreHandlerInterface,
	coverHandler hi.CoverHandlerInterface,
) *HandlersContainer {
	return &HandlersContainer{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		MakerHandler: makerHandler,
		CoreHandler:  coreHandler,
		CoverHandler: coverHandler,
	}
}