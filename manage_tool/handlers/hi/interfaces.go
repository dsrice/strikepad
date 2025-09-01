package hi

import "github.com/labstack/echo/v4"

// AuthHandlerInterface は認証ハンドラーのインターフェース
type AuthHandlerInterface interface {
	ShowLogin(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

// AdminHandlerInterface は管理者ハンドラーのインターフェース
type AdminHandlerInterface interface {
	ShowDashboard(c echo.Context) error
	ShowUsers(c echo.Context) error
	ShowUserDetail(c echo.Context) error
	UpdateUserStatus(c echo.Context) error
	DeleteUser(c echo.Context) error
}

// MakerHandlerInterface はメーカーハンドラーのインターフェース
type MakerHandlerInterface interface {
	ShowMakers(c echo.Context) error
	ShowMakerDetail(c echo.Context) error
	ShowCreateMaker(c echo.Context) error
	CreateMaker(c echo.Context) error
	ShowEditMaker(c echo.Context) error
	UpdateMaker(c echo.Context) error
	DeleteMaker(c echo.Context) error
	GetMakerStats(c echo.Context) error
	GetLogoPresignedURL(c echo.Context) error
	GetUploadPresignedURL(c echo.Context) error
	ConfirmLogoUpload(c echo.Context) error
}

// CoreHandlerInterface はコアハンドラーのインターフェース
type CoreHandlerInterface interface {
	ShowCores(c echo.Context) error
	ShowCreateCore(c echo.Context) error
	CreateCore(c echo.Context) error
	ShowCoreDetail(c echo.Context) error
	ShowEditCore(c echo.Context) error
	UpdateCore(c echo.Context) error
	DeleteCore(c echo.Context) error
	GetCoreStats(c echo.Context) error
}

// CoverHandlerInterface はカバーハンドラーのインターフェース
type CoverHandlerInterface interface {
	ShowCovers(c echo.Context) error
	ShowCreateCover(c echo.Context) error
	ShowEditCover(c echo.Context) error
	CreateCover(c echo.Context) error
	UpdateCover(c echo.Context) error
	ShowCoverDetail(c echo.Context) error
	DeleteCover(c echo.Context) error
	GetCoverStats(c echo.Context) error
}

// BallHandlerInterface はボールハンドラーのインターフェース
type BallHandlerInterface interface {
	ShowBalls(c echo.Context) error
	ShowCreateBall(c echo.Context) error
	CreateBall(c echo.Context) error
	ShowBallDetail(c echo.Context) error
	ShowEditBall(c echo.Context) error
	UpdateBall(c echo.Context) error
	DeleteBall(c echo.Context) error
	GetBallStats(c echo.Context) error
}