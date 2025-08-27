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

// CoreHandler はコア関連のHTTPハンドラー
type CoreHandler struct {
	coreRepo  *repository.CoreRepository
	makerRepo *repository.MakerRepository
}

// NewCoreHandler は新しいコアハンドラーを作成
func NewCoreHandler(coreRepo *repository.CoreRepository, makerRepo *repository.MakerRepository) *CoreHandler {
	return &CoreHandler{
		coreRepo:  coreRepo,
		makerRepo: makerRepo,
	}
}

// ShowCores はコア一覧画面を表示
func (h *CoreHandler) ShowCores(c echo.Context) error {
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

	var cores []*models.Core
	var totalCount int64

	// 検索条件に応じてコアを取得
	if search != "" {
		cores, err = h.coreRepo.Search(search, offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "コア検索エラー: "+err.Error())
		}
		// 検索時の総数は簡易的に現在のページの件数×10とする（実装簡化）
		totalCount = int64(len(cores) * 10)
	} else if makerIDParam != "" {
		makerID, parseErr := strconv.ParseUint(makerIDParam, 10, 32)
		if parseErr != nil {
			return c.String(http.StatusBadRequest, "無効なメーカーID")
		}
		cores, err = h.coreRepo.GetByMakerID(uint(makerID), offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "コア取得エラー: "+err.Error())
		}
		totalCount, err = h.coreRepo.CountByMaker(uint(makerID))
		if err != nil {
			totalCount = int64(len(cores))
		}
	} else {
		cores, err = h.coreRepo.GetAll(offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "コア取得エラー: "+err.Error())
		}
		totalCount, err = h.coreRepo.Count()
		if err != nil {
			totalCount = int64(len(cores))
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
		"Title":       "コア管理",
		"Cores":       cores,
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

	return c.Render(http.StatusOK, "cores.html", data)
}

// ShowCreateCore はコア作成画面を表示
func (h *CoreHandler) ShowCreateCore(c echo.Context) error {
	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":  "コア作成",
		"Makers": makers,
		"IsEdit": false,
	}

	return c.Render(http.StatusOK, "core_form.html", data)
}

// CreateCore は新しいコアを作成
func (h *CoreHandler) CreateCore(c echo.Context) error {
	// フォームデータを取得
	name := strings.TrimSpace(c.FormValue("name"))
	makerIDStr := strings.TrimSpace(c.FormValue("maker_id"))

	// バリデーション
	if name == "" {
		return c.String(http.StatusBadRequest, "コア名は必須です")
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

	// 新しいコアを作成（基本値を設定）
	core := &models.Core{
		Name:         name,
		MakerID:      uint(makerID),
		RG:           2.5,   // デフォルト値
		DeltaRG:      0.0,   // デフォルト値
		SymmetryFlag: false, // デフォルト値
	}

	err = h.coreRepo.Create(core)
	if err != nil {
		return c.String(http.StatusInternalServerError, "コア作成エラー: "+err.Error())
	}

	// コア一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/cores")
}

// ShowCoreDetail はコア詳細画面を表示
func (h *CoreHandler) ShowCoreDetail(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なコアIDです")
	}

	core, err := h.coreRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "コア取得エラー: "+err.Error())
	}
	if core == nil {
		return c.String(http.StatusNotFound, "コアが見つかりません")
	}

	data := map[string]interface{}{
		"Title": "コア詳細",
		"Core":  core,
	}

	return c.Render(http.StatusOK, "core_detail.html", data)
}

// DeleteCore はコアを論理削除
func (h *CoreHandler) DeleteCore(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なコアIDです"})
	}

	// コアが存在するかチェック
	core, err := h.coreRepo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "コア確認エラー"})
	}
	if core == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "コアが見つかりません"})
	}

	// 論理削除を実行
	err = h.coreRepo.Delete(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "コア削除エラー"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "コアが削除されました"})
}

// GetCoreStats はコア統計情報を取得
func (h *CoreHandler) GetCoreStats(c echo.Context) error {
	totalCount, err := h.coreRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "統計取得エラー"})
	}

	stats := map[string]interface{}{
		"total_cores": totalCount,
	}

	return c.JSON(http.StatusOK, stats)
}
