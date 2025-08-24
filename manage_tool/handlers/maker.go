package handlers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

const (
	jpgExt  = ".jpg"
	jpegExt = ".jpeg"
	pngExt  = ".png"
)

// MakerHandler はメーカー管理のハンドラー
type MakerHandler struct {
	makerRepo MakerRepositoryInterface
	s3Client  *config.S3Client
}

// NewMakerHandler は新しいメーカーハンドラーを作成
func NewMakerHandler(makerRepo MakerRepositoryInterface, s3Client *config.S3Client) *MakerHandler {
	return &MakerHandler{
		makerRepo: makerRepo,
		s3Client:  s3Client,
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
	log.Printf("uploadLogo開始: メーカーID=%d, ファイル名=%s", makerID, fileHeader.Filename)

	// ファイルサイズチェック (5MB)
	if fileHeader.Size > 5*1024*1024 {
		return fmt.Errorf("ファイルサイズが大きすぎます: %d bytes", fileHeader.Size)
	}

	// ファイル形式チェック
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != jpgExt && ext != jpegExt && ext != pngExt {
		return fmt.Errorf("サポートされていないファイル形式です: %s", ext)
	}
	log.Printf("ファイル検証OK: 拡張子=%s, サイズ=%d bytes", ext, fileHeader.Size)

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
		case jpgExt, jpegExt:
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		default:
			contentType = "application/octet-stream"
		}
	}

	log.Printf("S3アップロード開始: バケット=%s, キー=%s", h.s3Client.BucketName, objectName)
	_, err = h.s3Client.Client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket:        aws.String(h.s3Client.BucketName),
			Key:           aws.String(objectName),
			Body:          src,
			ContentLength: aws.Int64(fileHeader.Size),
			ContentType:   aws.String(contentType),
		},
	)
	if err != nil {
		return fmt.Errorf("S3へのアップロードに失敗しました: %w", err)
	}
	log.Printf("S3アップロード成功")

	// データベースにファイル名を保存
	log.Printf("データベース更新開始: メーカーID=%d, ファイル名=%s", makerID, fileHeader.Filename)
	err = h.makerRepo.UpdateLogoFile(makerID, fileHeader.Filename)
	if err != nil {
		return fmt.Errorf("データベースへのファイル名保存に失敗しました: %w", err)
	}
	log.Printf("データベース更新成功")

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
		duplicateMaker, duplicateErr := h.makerRepo.GetByName(name)
		if duplicateErr != nil {
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
		log.Printf("ロゴファイルが検出されました: %s, サイズ: %d bytes", file.Filename, file.Size)
		err = h.uploadLogo(uint(id), file)
		if err != nil {
			// ログ出力はするが、メーカー更新自体は成功とする
			log.Printf("ロゴアップロードに失敗しました: %v", err)
		} else {
			log.Printf("ロゴアップロードが成功しました: %s", file.Filename)
		}
	} else if err != nil {
		log.Printf("FormFileの取得でエラー: %v", err)
	} else {
		log.Printf("ロゴファイルが選択されていません")
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

	// 署名付きURLを生成（24時間有効）
	objectName := fmt.Sprintf("maker/%d/%s", makerID, maker.LogoFile)
	signedURL, err := h.s3Client.GeneratePresignedURL(objectName, 24*time.Hour)
	if err != nil {
		// エラーの場合は空文字を返す
		return ""
	}

	return signedURL
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

// GetLogoPresignedURL はメーカーのロゴ画像の署名付きURLを取得
func (h *MakerHandler) GetLogoPresignedURL(c echo.Context) error {
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

	// 署名付きURLを生成（1時間有効）
	objectName := fmt.Sprintf("maker/%d/%s", id, maker.LogoFile)
	signedURL, err := h.s3Client.GeneratePresignedURL(objectName, time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "署名付きURLの生成に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"url": signedURL,
	})
}

// GetUploadPresignedURL はロゴアップロード用の署名付きURLを取得
func (h *MakerHandler) GetUploadPresignedURL(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なメーカーIDです",
		})
	}

	// リクエストボディからファイル情報を取得
	var req struct {
		FileName    string `json:"fileName"`
		ContentType string `json:"contentType"`
	}

	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが無効です",
		})
	}

	// ファイル拡張子の検証
	ext := strings.ToLower(filepath.Ext(req.FileName))
	if ext != jpgExt && ext != jpegExt && ext != pngExt {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "サポートされていないファイル形式です",
		})
	}

	// アップロード用署名付きURLを生成（15分有効）
	objectName := fmt.Sprintf("maker/%d/%s", id, req.FileName)
	signedURL, err := h.s3Client.GenerateUploadPresignedURL(objectName, req.ContentType, 15*time.Minute)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "アップロード用署名付きURLの生成に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"uploadUrl": signedURL,
		"objectKey": objectName,
		"fileName":  req.FileName,
	})
}

// ConfirmLogoUpload はロゴアップロード完了をデータベースに記録
func (h *MakerHandler) ConfirmLogoUpload(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なメーカーIDです",
		})
	}

	// リクエストボディからファイル情報を取得
	var req struct {
		FileName string `json:"fileName"`
	}

	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが無効です",
		})
	}

	// データベースにファイル名を保存
	err = h.makerRepo.UpdateLogoFile(uint(id), req.FileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "データベースの更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ロゴファイルが正常にアップロードされました",
	})
}
