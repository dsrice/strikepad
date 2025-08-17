package handlers

import (
	"net/http"

	"strikepad-manage-tool/constants"

	"github.com/labstack/echo/v4"
)

// AuthHandler は認証関連のハンドラー
type AuthHandler struct {
	// 将来的にはデータベースやセッション管理を追加
}

// NewAuthHandler は新しい認証ハンドラーを作成
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// ShowLogin はログイン画面を表示
func (h *AuthHandler) ShowLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Error": "",
	})
}

// Login はログイン処理を実行
func (h *AuthHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 簡易認証（実際の実装では暗号化されたパスワードやデータベースを使用）
	if username == "admin" && password == "password" {
		// セッションを設定（簡易版）
		cookie := &http.Cookie{
			Name:  "session",
			Value: constants.SessionAuthenticated,
			Path:  "/",
		}
		c.SetCookie(cookie)

		return c.Redirect(http.StatusFound, "/dashboard")
	}

	// ログイン失敗
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Error": "ユーザー名またはパスワードが正しくありません",
	})
}

// ShowDashboard はダッシュボードを表示
func (h *AuthHandler) ShowDashboard(c echo.Context) error {
	// セッションチェック
	cookie, err := c.Cookie("session")
	if err != nil || cookie.Value != constants.SessionAuthenticated {
		return c.Redirect(http.StatusFound, "/login")
	}

	// ダッシュボードデータ（実際の実装ではデータベースから取得）
	data := map[string]interface{}{
		"Username": "admin",
		"Stats": map[string]interface{}{
			"TotalUsers":     150,
			"ActiveSessions": 23,
			"TodayLogins":    45,
		},
	}

	return c.Render(http.StatusOK, "dashboard.html", data)
}

// Logout はログアウト処理を実行
func (h *AuthHandler) Logout(c echo.Context) error {
	// クッキーを削除
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/login")
}
