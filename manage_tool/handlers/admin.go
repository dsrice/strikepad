package handlers

import (
	"net/http"
	"strconv"

	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository"

	"github.com/labstack/echo/v4"
)

// AdminHandler は管理機能のハンドラー
type AdminHandler struct {
	userRepo *repository.UserRepository
}

// NewAdminHandler は新しい管理ハンドラーを作成
func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

// ShowUsers はユーザー一覧画面を表示
func (h *AdminHandler) ShowUsers(c echo.Context) error {
	// セッションチェック
	cookie, err := c.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		return c.Redirect(http.StatusFound, "/login")
	}

	// ページネーション
	page, _ := strconv.Atoi(c.QueryParam("page"))
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

	totalCount, _ = h.userRepo.GetTotalUsersCount()

	data := map[string]interface{}{
		"Username":   "admin",
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
func (h *AdminHandler) ShowUserDetail(c echo.Context) error {
	// セッションチェック
	cookie, err := c.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
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
		"Username": "admin",
		"User":     user,
	}

	return c.Render(http.StatusOK, "user_detail.html", data)
}

// UpdateUserStatus はユーザーのアクティブ状態を更新
func (h *AdminHandler) UpdateUserStatus(c echo.Context) error {
	// セッションチェック
	cookie, err := c.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
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
func (h *AdminHandler) DeleteUser(c echo.Context) error {
	// セッションチェック
	cookie, err := c.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
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
