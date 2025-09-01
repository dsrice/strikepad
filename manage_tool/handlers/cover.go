package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"strikepad-manage-tool/handlers/hi"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository/ri"

	"github.com/labstack/echo/v4"
)

// coverHandler はカバー関連のHTTPハンドラー
type coverHandler struct {
	coverRepo CoverRepositoryInterface
	makerRepo MakerRepositoryInterface
}

// NewCoverHandlerInterface はDI用のCoverHandlerInterfaceを返す
func NewCoverHandlerInterface(coverRepo ri.CoverRepositoryInterface, makerRepo ri.MakerRepositoryInterface) hi.CoverHandlerInterface {
	return &coverHandler{
		coverRepo: coverRepo,
		makerRepo: makerRepo,
	}
}

// ShowCovers はカバー一覧画面を表示
func (h *coverHandler) ShowCovers(c echo.Context) error {
	// 検索フォームをバインド
	searchForm := new(SearchForm)
	if err := c.Bind(searchForm); err != nil {
		return c.String(http.StatusBadRequest, "検索パラメータが無効です")
	}

	page := searchForm.GetValidPage()
	search := strings.TrimSpace(searchForm.Search)
	makerIDParam := searchForm.MakerID

	const itemsPerPage = 10
	offset := (page - 1) * itemsPerPage

	var covers []*models.Cover
	var totalCount int64
	var err error

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
	makers, err := h.makerRepo.GetAllSimple()
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
		"CurrentPage": "covers",
		"Covers":      covers,
		"Makers":      makers,
		"Page":        page,
		"TotalPages":  totalPages,
		"TotalCount":  totalCount,
		"Search":      search,
		"MakerID":     makerIDParam,
		"HasPrevPage": page > 1,
		"HasNextPage": page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	return c.Render(http.StatusOK, "covers_new.html", data)
}

// ShowCreateCover はカバー作成画面を表示
func (h *coverHandler) ShowCreateCover(c echo.Context) error {
	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":       "カバー作成",
		"CurrentPage": "covers",
		"Makers":      makers,
		"IsEdit":      false,
	}

	return c.Render(http.StatusOK, "cover_form.html", data)
}

// ShowEditCover はカバー編集画面を表示
func (h *coverHandler) ShowEditCover(c echo.Context) error {
	// カバーIDを取得
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なカバーIDです")
	}

	// カバー情報を取得
	cover, err := h.coverRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー取得エラー: "+err.Error())
	}
	if cover == nil {
		return c.String(http.StatusNotFound, "カバーが見つかりません")
	}

	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":       "カバー編集",
		"CurrentPage": "covers",
		"Cover":       cover,
		"Makers":      makers,
		"IsEdit":      true,
	}

	return c.Render(http.StatusOK, "cover_form.html", data)
}

// CreateCover は新しいカバーを作成
func (h *coverHandler) CreateCover(c echo.Context) error {
	// フォームデータをバインド
	form := new(CoverForm)
	if bindErr := c.Bind(form); bindErr != nil {
		return c.String(http.StatusBadRequest, "フォームデータが無効です")
	}

	// バリデーション
	if validateErr := form.Validate(); validateErr != nil {
		return c.String(http.StatusBadRequest, validateErr.Error())
	}

	// メーカーが存在するかチェック
	maker, err := h.makerRepo.GetByID(form.MakerID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー確認エラー: "+err.Error())
	}
	if maker == nil {
		return c.String(http.StatusBadRequest, "指定されたメーカーが見つかりません")
	}

	// 新しいカバーを作成
	cover := &models.Cover{
		Name:         form.Name,
		MakerID:      form.MakerID,
		MaterialType: form.MaterialType,
		Rank:         form.Rank,
		IsDeleted:    false,
	}

	err = h.coverRepo.Create(cover)
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー作成エラー: "+err.Error())
	}

	// カバー一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/covers")
}

// UpdateCover はカバー情報を更新
func (h *coverHandler) UpdateCover(c echo.Context) error {
	// カバーIDを取得
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なカバーIDです")
	}

	// 既存のカバーが存在するかチェック
	existingCover, err := h.coverRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー確認エラー: "+err.Error())
	}
	if existingCover == nil {
		return c.String(http.StatusNotFound, "カバーが見つかりません")
	}

	// フォームデータをバインド
	form := new(CoverForm)
	if bindErr := c.Bind(form); bindErr != nil {
		return c.String(http.StatusBadRequest, "フォームデータが無効です")
	}

	// バリデーション
	if validateErr := form.Validate(); validateErr != nil {
		return c.String(http.StatusBadRequest, validateErr.Error())
	}

	// メーカーが存在するかチェック
	maker, makerErr := h.makerRepo.GetByID(form.MakerID)
	if makerErr != nil {
		return c.String(http.StatusInternalServerError, "メーカー確認エラー: "+makerErr.Error())
	}
	if maker == nil {
		return c.String(http.StatusBadRequest, "指定されたメーカーが見つかりません")
	}

	// カバー情報を更新
	existingCover.Name = form.Name
	existingCover.MakerID = form.MakerID
	existingCover.MaterialType = form.MaterialType
	existingCover.Rank = form.Rank

	err = h.coverRepo.Update(existingCover)
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー更新エラー: "+err.Error())
	}

	// カバー一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/covers")
}

// ShowCoverDetail はカバー詳細画面を表示
func (h *coverHandler) ShowCoverDetail(c echo.Context) error {
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
func (h *coverHandler) DeleteCover(c echo.Context) error {
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
func (h *coverHandler) GetCoverStats(c echo.Context) error {
	totalCount, err := h.coverRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "統計取得エラー"})
	}

	stats := map[string]interface{}{
		"total_covers": totalCount,
	}

	return c.JSON(http.StatusOK, stats)
}
