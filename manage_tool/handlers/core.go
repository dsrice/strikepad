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

// coreHandler はコア関連のHTTPハンドラー
type coreHandler struct {
	coreRepo  CoreRepositoryInterface
	makerRepo MakerRepositoryInterface
}

// NewCoreHandlerInterface はDI用のCoreHandlerInterfaceを返す
func NewCoreHandlerInterface(coreRepo ri.CoreRepositoryInterface, makerRepo ri.MakerRepositoryInterface) hi.CoreHandlerInterface {
	return &coreHandler{
		coreRepo:  coreRepo,
		makerRepo: makerRepo,
	}
}

// ShowCores はコア一覧画面を表示
func (h *coreHandler) ShowCores(c echo.Context) error {
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
func (h *coreHandler) ShowCreateCore(c echo.Context) error {
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
func (h *coreHandler) CreateCore(c echo.Context) error {
	// フォームデータを取得
	name := strings.TrimSpace(c.FormValue("name"))
	makerIDStr := strings.TrimSpace(c.FormValue("maker_id"))
	rgStr := strings.TrimSpace(c.FormValue("rg"))
	deltaRgStr := strings.TrimSpace(c.FormValue("delta_rg"))
	initDiffStr := strings.TrimSpace(c.FormValue("init_diff"))
	symmetryFlag := c.FormValue("symmetry_flag") == "1"

	// バリデーション
	if name == "" {
		return h.showFormWithError(c, "コア名は必須です", nil, false)
	}

	makerID, err := strconv.ParseUint(makerIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なメーカーIDです", nil, false)
	}

	// フォームデータをパース
	rg := float32(2.5) // デフォルト値
	if rgStr != "" {
		rgFloat, err := strconv.ParseFloat(rgStr, 32)
		if err != nil {
			return h.showFormWithError(c, "RGの値が無効です", nil, false)
		}
		rg = float32(rgFloat)
	}

	deltaRg := float32(0.0) // デフォルト値
	if deltaRgStr != "" {
		deltaRgFloat, err := strconv.ParseFloat(deltaRgStr, 32)
		if err != nil {
			return h.showFormWithError(c, "ΔRGの値が無効です", nil, false)
		}
		deltaRg = float32(deltaRgFloat)
	}

	var initDiff *float32
	if initDiffStr != "" {
		initDiffFloat, err := strconv.ParseFloat(initDiffStr, 32)
		if err != nil {
			return h.showFormWithError(c, "InitDiffの値が無効です", nil, false)
		}
		initDiffFloat32 := float32(initDiffFloat)
		initDiff = &initDiffFloat32
	}

	// メーカーが存在するかチェック
	maker, err := h.makerRepo.GetByID(uint(makerID))
	if err != nil {
		return h.showFormWithError(c, "メーカー確認エラー: "+err.Error(), nil, false)
	}
	if maker == nil {
		return h.showFormWithError(c, "指定されたメーカーが見つかりません", nil, false)
	}

	// 新しいコアを作成
	core := &models.Core{
		Name:         name,
		MakerID:      uint(makerID),
		RG:           rg,
		DeltaRG:      deltaRg,
		InitDiff:     initDiff,
		SymmetryFlag: symmetryFlag,
	}

	err = h.coreRepo.Create(core)
	if err != nil {
		return h.showFormWithError(c, "コア作成エラー: "+err.Error(), nil, false)
	}

	// コア一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/cores")
}

// ShowCoreDetail はコア詳細画面を表示
func (h *coreHandler) ShowCoreDetail(c echo.Context) error {
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
func (h *coreHandler) DeleteCore(c echo.Context) error {
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
func (h *coreHandler) GetCoreStats(c echo.Context) error {
	totalCount, err := h.coreRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "統計取得エラー"})
	}

	stats := map[string]interface{}{
		"total_cores": totalCount,
	}

	return c.JSON(http.StatusOK, stats)
}

// ShowEditCore はコア編集画面を表示
func (h *coreHandler) ShowEditCore(c echo.Context) error {
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

	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":  "コア編集",
		"Core":   core,
		"Makers": makers,
		"IsEdit": true,
	}

	return c.Render(http.StatusOK, "core_form.html", data)
}

// UpdateCore はコア情報を更新
func (h *coreHandler) UpdateCore(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なコアIDです")
	}

	// 既存のコアを取得
	existingCore, err := h.coreRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "コア取得エラー: "+err.Error())
	}
	if existingCore == nil {
		return c.String(http.StatusNotFound, "コアが見つかりません")
	}

	// フォームデータを取得
	name := strings.TrimSpace(c.FormValue("name"))
	makerIDStr := strings.TrimSpace(c.FormValue("maker_id"))
	rgStr := strings.TrimSpace(c.FormValue("rg"))
	deltaRgStr := strings.TrimSpace(c.FormValue("delta_rg"))
	initDiffStr := strings.TrimSpace(c.FormValue("init_diff"))
	symmetryFlag := c.FormValue("symmetry_flag") == "1"

	// バリデーション
	if name == "" {
		return h.showFormWithError(c, "コア名は必須です", existingCore, true)
	}

	makerID, err := strconv.ParseUint(makerIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なメーカーIDです", existingCore, true)
	}

	// フォームデータをパース
	rg := float32(2.5) // デフォルト値
	if rgStr != "" {
		rgFloat, err := strconv.ParseFloat(rgStr, 32)
		if err != nil {
			return h.showFormWithError(c, "RGの値が無効です", existingCore, true)
		}
		rg = float32(rgFloat)
	}

	deltaRg := float32(0.0) // デフォルト値
	if deltaRgStr != "" {
		deltaRgFloat, err := strconv.ParseFloat(deltaRgStr, 32)
		if err != nil {
			return h.showFormWithError(c, "ΔRGの値が無効です", existingCore, true)
		}
		deltaRg = float32(deltaRgFloat)
	}

	var initDiff *float32
	if initDiffStr != "" {
		initDiffFloat, err := strconv.ParseFloat(initDiffStr, 32)
		if err != nil {
			return h.showFormWithError(c, "InitDiffの値が無効です", existingCore, true)
		}
		initDiffFloat32 := float32(initDiffFloat)
		initDiff = &initDiffFloat32
	}

	// メーカーが存在するかチェック
	maker, err := h.makerRepo.GetByID(uint(makerID))
	if err != nil {
		return h.showFormWithError(c, "メーカー確認エラー: "+err.Error(), existingCore, true)
	}
	if maker == nil {
		return h.showFormWithError(c, "指定されたメーカーが見つかりません", existingCore, true)
	}

	// コア情報を更新
	existingCore.Name = name
	existingCore.MakerID = uint(makerID)
	existingCore.RG = rg
	existingCore.DeltaRG = deltaRg
	existingCore.InitDiff = initDiff
	existingCore.SymmetryFlag = symmetryFlag

	err = h.coreRepo.Update(existingCore)
	if err != nil {
		return h.showFormWithError(c, "コア更新エラー: "+err.Error(), existingCore, true)
	}

	// コア一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/cores")
}

// showFormWithError はエラーメッセージ付きでフォームを表示するヘルパーメソッド
func (h *coreHandler) showFormWithError(c echo.Context, errorMsg string, core *models.Core, isEdit bool) error {
	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		makers = []*models.Maker{} // エラー時は空配列
	}

	title := "コア作成"
	if isEdit {
		title = "コア編集"
	}

	data := map[string]interface{}{
		"Title":  title,
		"Core":   core,
		"Makers": makers,
		"IsEdit": isEdit,
		"Error":  errorMsg,
	}

	return c.Render(http.StatusBadRequest, "core_form.html", data)
}
