package handlers

import (
	"net/http"
	"strconv"
	"time"

	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/utils"

	"github.com/labstack/echo/v4"
)

// MakerHandler はメーカー管理のハンドラー
type MakerHandler struct {
	makerRepo *repository.MakerRepository
}

// NewMakerHandler は新しいメーカーハンドラーを作成
func NewMakerHandler(makerRepo *repository.MakerRepository) *MakerHandler {
	return &MakerHandler{
		makerRepo: makerRepo,
	}
}

// ShowMakers はメーカー一覧画面を表示
func (h *MakerHandler) ShowMakers(c echo.Context) error {
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

	var makers []models.MakerListItem
	var totalCount int64

	if search != "" {
		makers, err = h.makerRepo.SearchMakers(search, offset, limit)
	} else {
		makers, err = h.makerRepo.GetAllMakers(offset, limit)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "メーカー取得に失敗しました",
		})
	}

	totalCount, err = h.makerRepo.GetTotalMakersCount()
	if err != nil {
		// エラーログ出力して継続（統計情報のエラーで全体が止まらないように）
		totalCount = 0
	}

	data := map[string]interface{}{
		"Username":   currentUser.Name,
		"LoginID":    currentUser.LoginID,
		"Makers":     makers,
		"Page":       page,
		"Search":     search,
		"TotalCount": totalCount,
		"HasNext":    int64(offset+limit) < totalCount,
		"HasPrev":    page > 1,
		"NextPage":   page + 1,
		"PrevPage":   page - 1,
	}

	return c.Render(http.StatusOK, "makers.html", data)
}

// ShowMakerDetail はメーカー詳細画面を表示
func (h *MakerHandler) ShowMakerDetail(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なメーカーIDです",
		})
	}

	maker, err := h.makerRepo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "メーカーが見つかりません",
		})
	}

	data := map[string]interface{}{
		"Username": currentUser.Name,
		"LoginID":  currentUser.LoginID,
		"Maker":    maker,
	}

	return c.Render(http.StatusOK, "maker_detail.html", data)
}

// ShowCreateMaker はメーカー作成画面を表示
func (h *MakerHandler) ShowCreateMaker(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	data := map[string]interface{}{
		"Username": currentUser.Name,
		"LoginID":  currentUser.LoginID,
		"Error":    "",
	}

	return c.Render(http.StatusOK, "maker_create.html", data)
}

// CreateMaker はメーカーを作成
func (h *MakerHandler) CreateMaker(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "認証が必要です",
		})
	}

	name := c.FormValue("name")
	if name == "" {
		return c.Render(http.StatusOK, "maker_create.html", map[string]interface{}{
			"Username": currentUser.Name,
			"LoginID":  currentUser.LoginID,
			"Error":    "メーカー名は必須です",
		})
	}

	// 同名チェック
	existingMaker, err := h.makerRepo.GetByName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "データベースエラーが発生しました",
		})
	}
	if existingMaker != nil {
		return c.Render(http.StatusOK, "maker_create.html", map[string]interface{}{
			"Username": currentUser.Name,
			"LoginID":  currentUser.LoginID,
			"Error":    "このメーカー名は既に存在します",
		})
	}

	// メーカー作成
	maker := &models.Maker{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	err = h.makerRepo.Create(maker)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "メーカーの作成に失敗しました",
		})
	}

	return c.Redirect(http.StatusFound, "/admin/makers")
}

// DeleteMaker はメーカーを削除
func (h *MakerHandler) DeleteMaker(c echo.Context) error {
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
			"error": "無効なメーカーIDです",
		})
	}

	err = h.makerRepo.Delete(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "メーカーの削除に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "メーカーを削除しました",
	})
}

// GetMakerStats はメーカー統計を取得
func (h *MakerHandler) GetMakerStats(c echo.Context) error {
	// セッションから管理者ユーザー情報を取得
	currentUser := utils.GetCurrentAdminUser(c)
	if currentUser == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "認証が必要です",
		})
	}

	stats, err := h.makerRepo.GetMakerStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "統計情報の取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, stats)
}