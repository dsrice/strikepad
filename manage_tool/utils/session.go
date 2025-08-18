package utils

import (
	"fmt"
	"time"

	"strikepad-manage-tool/middleware"
	"strikepad-manage-tool/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// SetAdminSession は管理者ユーザーのセッションを設定
func SetAdminSession(c echo.Context, adminUser *models.AdminUser) error {
	sess, err := session.Get(middleware.SessionName, c)
	if err != nil {
		return fmt.Errorf("session.Get failed: %w", err)
	}

	// セッションデータを設定
	sess.Values[middleware.SessionKeyUserID] = adminUser.ID
	sess.Values[middleware.SessionKeyUserName] = adminUser.Name
	sess.Values[middleware.SessionKeyLoginID] = adminUser.LoginID
	sess.Values[middleware.SessionKeyLastAccess] = time.Now().Unix()

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return fmt.Errorf("session.Save failed: %w", err)
	}

	return nil
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

	id, ok := userID.(uint)
	if !ok {
		return nil
	}

	name, ok := userName.(string)
	if !ok {
		return nil
	}

	login, ok := loginID.(string)
	if !ok {
		return nil
	}

	return &AdminUserSession{
		ID:      id,
		Name:    name,
		LoginID: login,
	}
}

// AdminUserSession はセッション用の管理者ユーザー情報
type AdminUserSession struct {
	Name    string `json:"name"`
	LoginID string `json:"login_id"`
	ID      uint   `json:"id"`
}
