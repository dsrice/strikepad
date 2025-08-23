package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository"

	"github.com/labstack/echo/v4"
)

// CoverHandler はカバー関連のHTTPハンドラー
type CoverHandler struct {
	coverRepo *repository.CoverRepository
	makerRepo *repository.MakerRepository
}

// NewCoverHandler は新しいカバーハンドラーを作成
func NewCoverHandler(coverRepo *repository.CoverRepository, makerRepo *repository.MakerRepository) *CoverHandler {
	return &CoverHandler{
		coverRepo: coverRepo,
		makerRepo: makerRepo,
	}
}

// ShowCovers はカバー一覧画面を表示
func (h *CoverHandler) ShowCovers(c echo.Context) error {
	// ページネーションパラメータを取得
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	if page <= 0 {
		page = 1
	}

	// 検索パラメータを取得
	search := strings.TrimSpace(c.QueryParam("search"))
	makerIDParam := strings.TrimSpace(c.QueryParam("maker_id"))

	const itemsPerPage = 10
	offset := (page - 1) * itemsPerPage

	var covers []*models.Cover
	var totalCount int64

	// 検索条件に応じてカバーを取得
	if search != "" {
		covers, err = h.coverRepo.Search(search, offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "カバー検索エラー: "+err.Error())
		}
		// 検索時の総数は簡易的に現在のページの件数×10とする（実装簡化）
		totalCount = int64(len(covers) * 10)
	} else if makerIDParam != "" {
		makerID, parseErr := strconv.ParseUint(makerIDParam, 10, 32)
		if parseErr != nil {
			return c.String(http.StatusBadRequest, "無効なメーカーID")
		}
		covers, err = h.coverRepo.GetByMakerID(uint(makerID), offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "カバー取得エラー: "+err.Error())
		}
		totalCount, err = h.coverRepo.CountByMaker(uint(makerID))
		if err != nil {
			totalCount = int64(len(covers))
		}
	} else {
		covers, err = h.coverRepo.GetAll(offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "カバー取得エラー: "+err.Error())
		}
		totalCount, err = h.coverRepo.Count()
		if err != nil {
			totalCount = int64(len(covers))
		}
	}

	// メーカー一覧を取得（フィルター用）
	makerList, err := h.makerRepo.GetAllMakers(0, 0) // 全件取得
	makers := make([]*models.Maker, len(makerList))
	for i, m := range makerList {
		makers[i] = &models.Maker{
			ID:   m.ID,
			Name: m.Name,
		}
	}
	if err != nil {
		makers = []*models.Maker{} // エラー時は空配列
	}

	// ページネーション情報を計算
	totalPages := int(math.Ceil(float64(totalCount) / float64(itemsPerPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	// テンプレートデータを準備
	data := map[string]interface{}{
		"Title":       "カバー管理",
		"Covers":      covers,
		"Makers":      makers,
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalCount":  totalCount,
		"Search":      search,
		"MakerID":     makerIDParam,
		"HasPrevPage": page > 1,
		"HasNextPage": page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	return c.Render(http.StatusOK, "covers.html", data)
}

// ShowCreateCover はカバー作成画面を表示
func (h *CoverHandler) ShowCreateCover(c echo.Context) error {
	// メーカー一覧を取得
	makerList, err := h.makerRepo.GetAllMakers(0, 0)
	makers := make([]*models.Maker, len(makerList))
	for i, m := range makerList {
		makers[i] = &models.Maker{
			ID:   m.ID,
			Name: m.Name,
		}
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":  "カバー作成",
		"Makers": makers,
		"IsEdit": false,
	}

	return c.Render(http.StatusOK, "cover_form.html", data)
}

// CreateCover は新しいカバーを作成
func (h *CoverHandler) CreateCover(c echo.Context) error {
	// フォームデータを取得
	name := strings.TrimSpace(c.FormValue("name"))
	makerIDStr := strings.TrimSpace(c.FormValue("maker_id"))

	// バリデーション
	if name == "" {
		return c.String(http.StatusBadRequest, "カバー名は必須です")
	}

	makerID, err := strconv.ParseUint(makerIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なメーカーIDです")
	}

	// メーカーが存在するかチェック
	maker, err := h.makerRepo.GetByID(uint(makerID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー確認エラー: "+err.Error())
	}
	if maker == nil {
		return c.String(http.StatusBadRequest, "指定されたメーカーが見つかりません")
	}

	// 新しいカバーを作成（基本値を設定）
	cover := &models.Cover{
		Name:         name,
		MakerID:      uint(makerID),
		MaterialType: 1, // デフォルト値
		Rank:         1, // デフォルト値
	}

	err = h.coverRepo.Create(cover)
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー作成エラー: "+err.Error())
	}

	// カバー一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/covers")
}

// ShowCoverDetail はカバー詳細画面を表示
func (h *CoverHandler) ShowCoverDetail(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なカバーIDです")
	}

	cover, err := h.coverRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー取得エラー: "+err.Error())
	}
	if cover == nil {
		return c.String(http.StatusNotFound, "カバーが見つかりません")
	}

	data := map[string]interface{}{
		"Title": "カバー詳細",
		"Cover": cover,
	}

	return c.Render(http.StatusOK, "cover_detail.html", data)
}

// DeleteCover はカバーを論理削除
func (h *CoverHandler) DeleteCover(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なカバーIDです"})
	}

	// カバーが存在するかチェック
	cover, err := h.coverRepo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カバー確認エラー"})
	}
	if cover == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "カバーが見つかりません"})
	}

	// 論理削除を実行
	err = h.coverRepo.Delete(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "カバー削除エラー"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "カバーが削除されました"})
}

// GetCoverStats はカバー統計情報を取得
func (h *CoverHandler) GetCoverStats(c echo.Context) error {
	totalCount, err := h.coverRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "統計取得エラー"})
	}

	stats := map[string]interface{}{
		"total_covers": totalCount,
	}

	return c.JSON(http.StatusOK, stats)
}
