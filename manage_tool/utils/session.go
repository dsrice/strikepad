package utils

import (
	"strikepad-manage-tool/middleware"
	"strikepad-manage-tool/models"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// SetAdminSession は管理者ユーザーのセッションを設定
func SetAdminSession(c echo.Context, adminUser *models.AdminUser) error {
	sess, err := session.Get(middleware.SessionName, c)
	if err != nil {
		return err
	}

	// セッションデータを設定
	sess.Values[middleware.SessionKeyUserID] = adminUser.ID
	sess.Values[middleware.SessionKeyUserName] = adminUser.Name
	sess.Values[middleware.SessionKeyLoginID] = adminUser.LoginID
	sess.Values[middleware.SessionKeyLastAccess] = time.Now()

	return sess.Save(c.Request(), c.Response())
}

// ClearAdminSession は管理者ユーザーのセッションをクリア
func ClearAdminSession(c echo.Context) error {
	sess, err := session.Get(middleware.SessionName, c)
	if err != nil {
		return err
	}

	// セッションデータをクリア
	sess.Values = make(map[interface{}]interface{})

	return sess.Save(c.Request(), c.Response())
}

// GetCurrentAdminUser は現在のセッションから管理者ユーザー情報を取得
func GetCurrentAdminUser(c echo.Context) *AdminUserSession {
	userID := c.Get("user_id")
	userName := c.Get("user_name")
	loginID := c.Get("login_id")

	if userID == nil {
		return nil
	}

	return &AdminUserSession{
		ID:      userID.(uint),
		Name:    userName.(string),
		LoginID: loginID.(string),
	}
}

// AdminUserSession はセッション用の管理者ユーザー情報
type AdminUserSession struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	LoginID string `json:"login_id"`
}