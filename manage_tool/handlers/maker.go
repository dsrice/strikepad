package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/utils"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

// MakerHandler はメーカー管理のハンドラー
type MakerHandler struct {
	makerRepo   MakerRepositoryInterface
	minioClient *config.MinIOClient
}

// NewMakerHandler は新しいメーカーハンドラーを作成
func NewMakerHandler(makerRepo MakerRepositoryInterface, minioClient *config.MinIOClient) *MakerHandler {
	return &MakerHandler{
		makerRepo:   makerRepo,
		minioClient: minioClient,
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

	// ロゴ画像のアップロード処理
	file, err := c.FormFile("logo")
	if err == nil && file != nil {
		err = h.uploadLogo(maker.ID, file)
		if err != nil {
			// ログ出力はするが、メーカー作成自体は成功とする
			fmt.Printf("ロゴアップロードに失敗しました: %v\n", err)
		}
	}

	return c.Redirect(http.StatusFound, "/admin/makers")
}

// uploadLogo はロゴ画像をMinIOにアップロードし、DBにファイル名を保存
func (h *MakerHandler) uploadLogo(makerID uint, fileHeader *multipart.FileHeader) error {
	// ファイルサイズチェック (5MB)
	if fileHeader.Size > 5*1024*1024 {
		return fmt.Errorf("ファイルサイズが大きすぎます: %d bytes", fileHeader.Size)
	}

	// ファイル形式チェック
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("サポートされていないファイル形式です: %s", ext)
	}

	// ファイルを開く
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer src.Close()

	// MinIOにアップロード
	objectName := fmt.Sprintf("maker/%d/%s", makerID, fileHeader.Filename)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		default:
			contentType = "application/octet-stream"
		}
	}

	_, err = h.minioClient.Client.PutObject(
		context.Background(),
		h.minioClient.BucketName,
		objectName,
		src,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return fmt.Errorf("MinIOへのアップロードに失敗しました: %w", err)
	}

	// データベースにファイル名を保存
	err = h.makerRepo.UpdateLogoFile(makerID, fileHeader.Filename)
	if err != nil {
		return fmt.Errorf("データベースへのファイル名保存に失敗しました: %w", err)
	}

	return nil
}

// ShowEditMaker はメーカー編集画面を表示
func (h *MakerHandler) ShowEditMaker(c echo.Context) error {
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

	// ロゴURLを生成
	logoURL := h.getLogoURL(uint(id))

	data := map[string]interface{}{
		"Username": currentUser.Name,
		"LoginID":  currentUser.LoginID,
		"Maker":    maker,
		"LogoURL":  logoURL,
		"Error":    "",
	}

	return c.Render(http.StatusOK, "maker_edit.html", data)
}

// UpdateMaker はメーカー情報を更新
func (h *MakerHandler) UpdateMaker(c echo.Context) error {
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

	// 既存のメーカー取得
	existingMaker, err := h.makerRepo.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "メーカーが見つかりません",
		})
	}

	name := c.FormValue("name")
	if name == "" {
		logoURL := h.getLogoURL(uint(id))
		return c.Render(http.StatusOK, "maker_edit.html", map[string]interface{}{
			"Username": currentUser.Name,
			"LoginID":  currentUser.LoginID,
			"Maker":    existingMaker,
			"LogoURL":  logoURL,
			"Error":    "メーカー名は必須です",
		})
	}

	// 同名チェック（自分以外で同じ名前がないか確認）
	if name != existingMaker.Name {
		duplicateMaker, err := h.makerRepo.GetByName(name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "データベースエラーが発生しました",
			})
		}
		if duplicateMaker != nil {
			logoURL := h.getLogoURL(uint(id))
			return c.Render(http.StatusOK, "maker_edit.html", map[string]interface{}{
				"Username": currentUser.Name,
				"LoginID":  currentUser.LoginID,
				"Maker":    existingMaker,
				"LogoURL":  logoURL,
				"Error":    "このメーカー名は既に存在します",
			})
		}
	}

	// メーカー情報を更新
	existingMaker.Name = name
	existingMaker.UpdatedAt = time.Now()

	err = h.makerRepo.Update(existingMaker)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "メーカーの更新に失敗しました",
		})
	}

	// ロゴ画像のアップロード処理（ファイルが選択されている場合のみ）
	file, err := c.FormFile("logo")
	if err == nil && file != nil {
		err = h.uploadLogo(uint(id), file)
		if err != nil {
			// ログ出力はするが、メーカー更新自体は成功とする
			fmt.Printf("ロゴアップロードに失敗しました: %v\n", err)
		}
	}

	return c.Redirect(http.StatusFound, "/admin/makers")
}

// getLogoURL はメーカーのロゴURLを取得
func (h *MakerHandler) getLogoURL(makerID uint) string {
	// データベースからロゴファイル名を取得
	maker, err := h.makerRepo.GetByID(makerID)
	if err != nil || maker == nil || maker.LogoFile == "" {
		return "" // ロゴファイルが登録されていない場合
	}

	// 管理ツール内のエンドポイントを使用
	return fmt.Sprintf("/admin/makers/%d/logo", makerID)
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

// ServeLogoImage はメーカーのロゴ画像を配信
func (h *MakerHandler) ServeLogoImage(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なメーカーIDです",
		})
	}

	// データベースからロゴファイル名を取得
	maker, err := h.makerRepo.GetByID(uint(id))
	if err != nil || maker == nil || maker.LogoFile == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "ロゴ画像が見つかりません",
		})
	}

	// MinIOからオブジェクトを取得
	objectName := fmt.Sprintf("maker/%d/%s", id, maker.LogoFile)
	object, err := h.minioClient.Client.GetObject(
		context.Background(),
		h.minioClient.BucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "画像ファイルが見つかりません",
		})
	}
	defer object.Close()

	// オブジェクト情報を取得
	objInfo, err := object.Stat()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "画像情報の取得に失敗しました",
		})
	}

	// ファイル拡張子からContent-Typeを決定
	ext := strings.ToLower(filepath.Ext(maker.LogoFile))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// レスポンスヘッダーを設定
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", objInfo.Size))
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")

	// オブジェクトの内容をレスポンスにストリーミング
	return c.Stream(http.StatusOK, contentType, object)
}