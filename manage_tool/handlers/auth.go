package handlers

import (
	"log"
	"net/http"
	"strings"

	"strikepad-manage-tool/utils"

	"github.com/labstack/echo/v4"
)

// AuthHandler は認証関連のハンドラー
type AuthHandler struct {
	adminRepo AdminRepositoryInterface
}

// NewAuthHandler は新しい認証ハンドラーを作成
func NewAuthHandler(adminRepo AdminRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		adminRepo: adminRepo,
	}
}

// ShowLogin はログイン画面を表示
func (h *AuthHandler) ShowLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Error": "",
	})
}

// Login はログイン処理を実行
func (h *AuthHandler) Login(c echo.Context) error {
	loginID := strings.TrimSpace(c.FormValue("login_id"))
	password := strings.TrimSpace(c.FormValue("password"))

	// 入力値検証
	if loginID == "" || password == "" {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "ログインIDとパスワードを入力してください",
		})
	}

	// データベースから管理者ユーザーを取得
	adminUser, err := h.adminRepo.GetByLoginID(loginID)
	if err != nil {
		log.Printf("Database error during login: %v", err)
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "システムエラーが発生しました",
		})
	}

	// ユーザーが存在しない場合
	if adminUser == nil {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "ログインIDまたはパスワードが正しくありません",
		})
	}

	// パスワード照合
	if !utils.CheckPassword(password, adminUser.Password) {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "ログインIDまたはパスワードが正しくありません",
		})
	}

	// セッションを設定
	log.Printf("Attempting to create session for user: %s (ID: %d)", adminUser.LoginID, adminUser.ID)
	err = utils.SetAdminSession(c, adminUser)
	if err != nil {
		log.Printf("Session error during login: %v", err)
		log.Printf("Session error type: %T", err)
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "セッションの作成に失敗しました",
		})
	}
	log.Printf("Session created successfully for user: %s", adminUser.LoginID)

	log.Printf("Admin user logged in: %s (%s)", adminUser.Name, adminUser.LoginID)
	return c.Redirect(http.StatusFound, "/dashboard")
}

// Logout はログアウト処理を実行
func (h *AuthHandler) Logout(c echo.Context) error {
	// セッションをクリア
	err := utils.ClearAdminSession(c)
	if err != nil {
		log.Printf("Session clear error during logout: %v", err)
	}

	// 現在のユーザー情報をログに記録
	if currentUser := utils.GetCurrentAdminUser(c); currentUser != nil {
		log.Printf("Admin user logged out: %s (%s)", currentUser.Name, currentUser.LoginID)
	}

	return c.Redirect(http.StatusFound, "/login")
}
