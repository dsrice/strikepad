package utils

import (
	"strikepad-manage-tool/models"

	"github.com/labstack/echo/v4"
)

// テスト用のセッション設定ヘルパー関数
func SetCurrentAdminUser(c echo.Context, user *models.AdminUser) {
	// テスト環境では直接コンテキストに設定
	c.Set("user_id", user.ID)
	c.Set("user_name", user.Name)
	c.Set("login_id", user.LoginID)
}
