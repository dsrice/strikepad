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

// 型エイリアス
type BallRepositoryInterface = ri.BallRepositoryInterface

// ballHandler はボール関連のHTTPハンドラー
type ballHandler struct {
	ballRepo  BallRepositoryInterface
	makerRepo MakerRepositoryInterface
	coreRepo  CoreRepositoryInterface
	coverRepo CoverRepositoryInterface
}

// NewBallHandlerInterface はDI用のBallHandlerInterfaceを返す
func NewBallHandlerInterface(
	ballRepo ri.BallRepositoryInterface,
	makerRepo ri.MakerRepositoryInterface,
	coreRepo ri.CoreRepositoryInterface,
	coverRepo ri.CoverRepositoryInterface,
) hi.BallHandlerInterface {
	return &ballHandler{
		ballRepo:  ballRepo,
		makerRepo: makerRepo,
		coreRepo:  coreRepo,
		coverRepo: coverRepo,
	}
}

// ShowBalls はボール一覧画面を表示
func (h *ballHandler) ShowBalls(c echo.Context) error {
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

	var balls []*models.Ball
	var totalCount int64

	// 検索条件に応じてボールを取得
	if search != "" {
		balls, err = h.ballRepo.Search(search, offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "ボール検索エラー: "+err.Error())
		}
		// 検索時の総数は簡易的に現在のページの件数×10とする（実装簡化）
		totalCount = int64(len(balls) * 10)
	} else if makerIDParam != "" {
		makerID, parseErr := strconv.ParseUint(makerIDParam, 10, 32)
		if parseErr != nil {
			return c.String(http.StatusBadRequest, "無効なメーカーID")
		}
		balls, err = h.ballRepo.GetByMakerID(uint(makerID), offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "ボール取得エラー: "+err.Error())
		}
		totalCount, err = h.ballRepo.CountByMaker(uint(makerID))
		if err != nil {
			totalCount = int64(len(balls))
		}
	} else {
		balls, err = h.ballRepo.GetAll(offset, itemsPerPage)
		if err != nil {
			return c.String(http.StatusInternalServerError, "ボール取得エラー: "+err.Error())
		}
		totalCount, err = h.ballRepo.Count()
		if err != nil {
			totalCount = int64(len(balls))
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
		"Title":       "ボール管理",
		"Balls":       balls,
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

	return c.Render(http.StatusOK, "balls.html", data)
}

// ShowCreateBall はボール作成画面を表示
func (h *ballHandler) ShowCreateBall(c echo.Context) error {
	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	// コア一覧を取得（簡易版）
	cores, err := h.coreRepo.GetAll(0, 1000) // 最大1000件
	if err != nil {
		return c.String(http.StatusInternalServerError, "コア取得エラー: "+err.Error())
	}

	// カバー一覧を取得（簡易版）
	covers, err := h.coverRepo.GetAll(0, 1000) // 最大1000件
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":  "ボール作成",
		"Makers": makers,
		"Cores":  cores,
		"Covers": covers,
		"IsEdit": false,
	}

	return c.Render(http.StatusOK, "ball_form.html", data)
}

// CreateBall は新しいボールを作成
func (h *ballHandler) CreateBall(c echo.Context) error {
	// フォームデータを取得
	name := strings.TrimSpace(c.FormValue("name"))
	url := strings.TrimSpace(c.FormValue("url"))
	makerIDStr := strings.TrimSpace(c.FormValue("maker_id"))
	coreIDStr := strings.TrimSpace(c.FormValue("core_id"))
	coverIDStr := strings.TrimSpace(c.FormValue("cover_id"))

	// バリデーション
	if name == "" {
		return h.showFormWithError(c, "ボール名は必須です", nil, false)
	}
	if url == "" {
		return h.showFormWithError(c, "URLは必須です", nil, false)
	}

	makerID, err := strconv.ParseUint(makerIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なメーカーIDです", nil, false)
	}

	coreID, err := strconv.ParseUint(coreIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なコアIDです", nil, false)
	}

	coverID, err := strconv.ParseUint(coverIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なカバーIDです", nil, false)
	}

	// 関連データが存在するかチェック
	maker, err := h.makerRepo.GetByID(uint(makerID))
	if err != nil {
		return h.showFormWithError(c, "メーカー確認エラー: "+err.Error(), nil, false)
	}
	if maker == nil {
		return h.showFormWithError(c, "指定されたメーカーが見つかりません", nil, false)
	}

	core, err := h.coreRepo.GetByID(uint(coreID))
	if err != nil {
		return h.showFormWithError(c, "コア確認エラー: "+err.Error(), nil, false)
	}
	if core == nil {
		return h.showFormWithError(c, "指定されたコアが見つかりません", nil, false)
	}

	cover, err := h.coverRepo.GetByID(uint(coverID))
	if err != nil {
		return h.showFormWithError(c, "カバー確認エラー: "+err.Error(), nil, false)
	}
	if cover == nil {
		return h.showFormWithError(c, "指定されたカバーが見つかりません", nil, false)
	}

	// 新しいボールを作成
	ball := &models.Ball{
		Name:    name,
		URL:     url,
		MakerID: uint(makerID),
		CoreID:  uint(coreID),
		CoverID: uint(coverID),
	}

	err = h.ballRepo.Create(ball)
	if err != nil {
		return h.showFormWithError(c, "ボール作成エラー: "+err.Error(), nil, false)
	}

	// ボール一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/balls")
}

// ShowBallDetail はボール詳細画面を表示
func (h *ballHandler) ShowBallDetail(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なボールIDです")
	}

	ball, err := h.ballRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "ボール取得エラー: "+err.Error())
	}
	if ball == nil {
		return c.String(http.StatusNotFound, "ボールが見つかりません")
	}

	data := map[string]interface{}{
		"Title": "ボール詳細",
		"Ball":  ball,
	}

	return c.Render(http.StatusOK, "ball_detail.html", data)
}

// ShowEditBall はボール編集画面を表示
func (h *ballHandler) ShowEditBall(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なボールIDです")
	}

	ball, err := h.ballRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "ボール取得エラー: "+err.Error())
	}
	if ball == nil {
		return c.String(http.StatusNotFound, "ボールが見つかりません")
	}

	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		return c.String(http.StatusInternalServerError, "メーカー取得エラー: "+err.Error())
	}

	// コア一覧を取得
	cores, err := h.coreRepo.GetAll(0, 1000) // 最大1000件
	if err != nil {
		return c.String(http.StatusInternalServerError, "コア取得エラー: "+err.Error())
	}

	// カバー一覧を取得
	covers, err := h.coverRepo.GetAll(0, 1000) // 最大1000件
	if err != nil {
		return c.String(http.StatusInternalServerError, "カバー取得エラー: "+err.Error())
	}

	data := map[string]interface{}{
		"Title":  "ボール編集",
		"Ball":   ball,
		"Makers": makers,
		"Cores":  cores,
		"Covers": covers,
		"IsEdit": true,
	}

	return c.Render(http.StatusOK, "ball_form.html", data)
}

// UpdateBall はボール情報を更新
func (h *ballHandler) UpdateBall(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "無効なボールIDです")
	}

	// 既存のボールを取得
	existingBall, err := h.ballRepo.GetByID(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "ボール取得エラー: "+err.Error())
	}
	if existingBall == nil {
		return c.String(http.StatusNotFound, "ボールが見つかりません")
	}

	// フォームデータを取得
	name := strings.TrimSpace(c.FormValue("name"))
	url := strings.TrimSpace(c.FormValue("url"))
	makerIDStr := strings.TrimSpace(c.FormValue("maker_id"))
	coreIDStr := strings.TrimSpace(c.FormValue("core_id"))
	coverIDStr := strings.TrimSpace(c.FormValue("cover_id"))

	// バリデーション
	if name == "" {
		return h.showFormWithError(c, "ボール名は必須です", existingBall, true)
	}
	if url == "" {
		return h.showFormWithError(c, "URLは必須です", existingBall, true)
	}

	makerID, err := strconv.ParseUint(makerIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なメーカーIDです", existingBall, true)
	}

	coreID, err := strconv.ParseUint(coreIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なコアIDです", existingBall, true)
	}

	coverID, err := strconv.ParseUint(coverIDStr, 10, 32)
	if err != nil {
		return h.showFormWithError(c, "無効なカバーIDです", existingBall, true)
	}

	// 関連データが存在するかチェック
	maker, err := h.makerRepo.GetByID(uint(makerID))
	if err != nil {
		return h.showFormWithError(c, "メーカー確認エラー: "+err.Error(), existingBall, true)
	}
	if maker == nil {
		return h.showFormWithError(c, "指定されたメーカーが見つかりません", existingBall, true)
	}

	core, err := h.coreRepo.GetByID(uint(coreID))
	if err != nil {
		return h.showFormWithError(c, "コア確認エラー: "+err.Error(), existingBall, true)
	}
	if core == nil {
		return h.showFormWithError(c, "指定されたコアが見つかりません", existingBall, true)
	}

	cover, err := h.coverRepo.GetByID(uint(coverID))
	if err != nil {
		return h.showFormWithError(c, "カバー確認エラー: "+err.Error(), existingBall, true)
	}
	if cover == nil {
		return h.showFormWithError(c, "指定されたカバーが見つかりません", existingBall, true)
	}

	// ボール情報を更新
	existingBall.Name = name
	existingBall.URL = url
	existingBall.MakerID = uint(makerID)
	existingBall.CoreID = uint(coreID)
	existingBall.CoverID = uint(coverID)

	err = h.ballRepo.Update(existingBall)
	if err != nil {
		return h.showFormWithError(c, "ボール更新エラー: "+err.Error(), existingBall, true)
	}

	// ボール一覧にリダイレクト
	return c.Redirect(http.StatusSeeOther, "/admin/balls")
}

// DeleteBall はボールを論理削除
func (h *ballHandler) DeleteBall(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なボールIDです"})
	}

	// ボールが存在するかチェック
	ball, err := h.ballRepo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ボール確認エラー"})
	}
	if ball == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ボールが見つかりません"})
	}

	// 論理削除を実行
	err = h.ballRepo.Delete(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ボール削除エラー"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ボールが削除されました"})
}

// GetBallStats はボール統計情報を取得
func (h *ballHandler) GetBallStats(c echo.Context) error {
	totalCount, err := h.ballRepo.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "統計取得エラー"})
	}

	stats := map[string]interface{}{
		"total_balls": totalCount,
	}

	return c.JSON(http.StatusOK, stats)
}

// showFormWithError はエラーメッセージ付きでフォームを表示するヘルパーメソッド
func (h *ballHandler) showFormWithError(c echo.Context, errorMsg string, ball *models.Ball, isEdit bool) error {
	// メーカー一覧を取得
	makers, err := h.makerRepo.GetAllSimple()
	if err != nil {
		makers = []*models.Maker{} // エラー時は空配列
	}

	// コア一覧を取得
	cores, err := h.coreRepo.GetAll(0, 1000) // 最大1000件
	if err != nil {
		cores = []*models.Core{} // エラー時は空配列
	}

	// カバー一覧を取得
	covers, err := h.coverRepo.GetAll(0, 1000) // 最大1000件
	if err != nil {
		covers = []*models.Cover{} // エラー時は空配列
	}

	title := "ボール作成"
	if isEdit {
		title = "ボール編集"
	}

	data := map[string]interface{}{
		"Title":  title,
		"Ball":   ball,
		"Makers": makers,
		"Cores":  cores,
		"Covers": covers,
		"IsEdit": isEdit,
		"Error":  errorMsg,
	}

	return c.Render(http.StatusBadRequest, "ball_form.html", data)
}