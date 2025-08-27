package handlers

import (
	"net/http"
	"strconv"

	"strikepad-manage-tool/handlers/hi"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository/ri"
	"strikepad-manage-tool/utils"

	"github.com/labstack/echo/v4"
)

// adminHandler は管理機能のハンドラー
type adminHandler struct {
	userRepo UserRepositoryInterface
}

// NewAdminHandlerInterface はDI用のAdminHandlerInterfaceを返す
func NewAdminHandlerInterface(userRepo ri.UserRepositoryInterface) hi.AdminHandlerInterface {
	return &adminHandler{
		userRepo: userRepo,
	}
}

// ShowDashboard はダッシュボードを表示
func (h *adminHandler) ShowDashboard(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	// 統計情報を取得
	stats, err := h.userRepo.GetUserStats()
	if err != nil {
		// エラーが発生した場合はデフォルト値を設定
		stats = &models.UserStats{
			TotalUsers:     0,
			ActiveUsers:    0,
			InactiveUsers:  0,
			TodayLogins:    0,
			ActiveSessions: 0,
		}
	}

	data := map[string]interface{}{
		"Username": currentUser.Name,
		"LoginID":  currentUser.LoginID,
		"Stats":    stats,
	}

	return c.Render(http.StatusOK, "dashboard.html", data)
}

// ShowUsers はユーザー一覧画面を表示
func (h *adminHandler) ShowUsers(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	// ページネーション
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	// 検索クエリ
	search := c.QueryParam("search")

	var users []models.UserListItem
	var totalCount int64

	if search != "" {
		users, err = h.userRepo.SearchUsers(search, offset, limit)
	} else {
		users, err = h.userRepo.GetAllUsers(offset, limit)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ユーザー取得に失敗しました",
		})
	}

	totalCount, err = h.userRepo.GetTotalUsersCount()
	if err != nil {
		// エラーログ出力して継続（統計情報のエラーで全体が止まらないように）
		totalCount = 0
	}

	data := map[string]interface{}{
		"Username":   currentUser.Name,
		"LoginID":    currentUser.LoginID,
		"Users":      users,
		"Page":       page,
		"Search":     search,
		"TotalCount": totalCount,
		"HasNext":    int64(offset+limit) < totalCount,
		"HasPrev":    page > 1,
		"NextPage":   page + 1,
		"PrevPage":   page - 1,
	}

	return c.Render(http.StatusOK, "users.html", data)
}

// ShowUserDetail はユーザー詳細画面を表示
func (h *adminHandler) ShowUserDetail(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なユーザーIDです",
		})
	}

	user, err := h.userRepo.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "ユーザーが見つかりません",
		})
	}

	data := map[string]interface{}{
		"Username": currentUser.Name,
		"LoginID":  currentUser.LoginID,
		"User":     user,
	}

	return c.Render(http.StatusOK, "user_detail.html", data)
}

// UpdateUserStatus はユーザーのアクティブ状態を更新
func (h *adminHandler) UpdateUserStatus(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "認証が必要です",
		})
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なユーザーIDです",
		})
	}

	isActive := c.FormValue("is_active") == "true"

	err = h.userRepo.UpdateUserStatus(uint(id), isActive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ユーザー状態の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ユーザー状態を更新しました",
	})
}

// DeleteUser はユーザーを削除
func (h *adminHandler) DeleteUser(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "認証が必要です",
		})
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なユーザーIDです",
		})
	}

	err = h.userRepo.DeleteUser(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ユーザーの削除に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ユーザーを削除しました",
	})
}
