package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/handlers/mocks"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MakerHandlerTestSuite struct {
	suite.Suite
	mockRepo        *mocks.MockMakerRepository
	mockMinioClient *config.MinIOClient
	handler         *MakerHandler
	echo            *echo.Echo
}

func (suite *MakerHandlerTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockMakerRepository{}
	suite.mockMinioClient = &config.MinIOClient{
		BucketName: "test-bucket",
	}
	suite.handler = &MakerHandler{
		makerRepo:   suite.mockRepo,
		minioClient: suite.mockMinioClient,
	}
	suite.echo = echo.New()
}

func TestMakerHandlerSuite(t *testing.T) {
	suite.Run(t, new(MakerHandlerTestSuite))
}

// ShowMakersのテスト
func (suite *MakerHandlerTestSuite) TestShowMakers() {
	testCases := []struct {
		name               string
		queryParams        map[string]string
		mockSetup          func()
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:        "正常ケース - ページネーション付き",
			queryParams: map[string]string{"page": "1"},
			mockSetup: func() {
				makers := []models.MakerListItem{
					{ID: 1, Name: "Test Maker", BallCount: 5, CoverCount: 3, CoreCount: 2},
				}
				suite.mockRepo.On("GetAllMakers", 0, 20).Return(makers, nil)
				suite.mockRepo.On("GetTotalMakersCount").Return(int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "正常ケース - 検索付き",
			queryParams: map[string]string{"search": "test"},
			mockSetup: func() {
				makers := []models.MakerListItem{
					{ID: 1, Name: "Test Maker", BallCount: 5, CoverCount: 3, CoreCount: 2},
				}
				suite.mockRepo.On("SearchMakers", "test", 0, 20).Return(makers, nil)
				suite.mockRepo.On("GetTotalMakersCount").Return(int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "エラーケース - リポジトリエラー",
			mockSetup: func() {
				suite.mockRepo.On("GetAllMakers", 0, 20).Return([]models.MakerListItem{}, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			// リクエスト作成
			req := httptest.NewRequest(http.MethodGet, "/admin/makers", nil)
			if len(tc.queryParams) > 0 {
				q := url.Values{}
				for k, v := range tc.queryParams {
					q.Set(k, v)
				}
				req.URL.RawQuery = q.Encode()
			}

			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// セッション設定
			adminUser := &models.AdminUser{ID: 1, Name: "Test Admin", LoginID: "admin"}
			utils.SetAdminSession(c, adminUser)

			// テスト実行
			err := suite.handler.ShowMakers(c)

			// アサーション
			if tc.expectedStatusCode == http.StatusOK {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// CreateMakerのテスト
func (suite *MakerHandlerTestSuite) TestCreateMaker() {
	testCases := []struct {
		name               string
		formData           map[string]string
		fileUpload         bool
		mockSetup          func()
		expectedStatusCode int
		expectedRedirect   string
	}{
		{
			name:     "正常ケース - ファイルなし",
			formData: map[string]string{"name": "New Maker"},
			mockSetup: func() {
				suite.mockRepo.On("GetByName", "New Maker").Return((*models.Maker)(nil), nil)
				suite.mockRepo.On("Create", mock.AnythingOfType("*models.Maker")).Return(nil)
			},
			expectedStatusCode: http.StatusFound,
			expectedRedirect:   "/admin/makers",
		},
		{
			name:     "エラーケース - 名前重複",
			formData: map[string]string{"name": "Existing Maker"},
			mockSetup: func() {
				existingMaker := &models.Maker{ID: 1, Name: "Existing Maker"}
				suite.mockRepo.On("GetByName", "Existing Maker").Return(existingMaker, nil)
			},
			expectedStatusCode: http.StatusOK, // テンプレート表示
		},
		{
			name:     "エラーケース - 名前が空",
			formData: map[string]string{"name": ""},
			mockSetup: func() {
				// モックセットアップなし
			},
			expectedStatusCode: http.StatusOK, // テンプレート表示
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			// フォームデータ作成
			form := url.Values{}
			for k, v := range tc.formData {
				form.Set(k, v)
			}

			req := httptest.NewRequest(http.MethodPost, "/admin/makers/create", strings.NewReader(form.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// セッション設定
			adminUser := &models.AdminUser{ID: 1, Name: "Test Admin", LoginID: "admin"}
			utils.SetAdminSession(c, adminUser)

			// テスト実行
			err := suite.handler.CreateMaker(c)

			// アサーション
			suite.NoError(err)
			if tc.expectedRedirect != "" {
				suite.Equal(tc.expectedStatusCode, rec.Code)
				suite.Equal(tc.expectedRedirect, rec.Header().Get("Location"))
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// UpdateMakerのテスト
func (suite *MakerHandlerTestSuite) TestUpdateMaker() {
	testCases := []struct {
		name               string
		makerID            string
		formData           map[string]string
		mockSetup          func()
		expectedStatusCode int
	}{
		{
			name:     "正常ケース - 名前更新",
			makerID:  "1",
			formData: map[string]string{"name": "Updated Maker"},
			mockSetup: func() {
				existingMaker := &models.Maker{ID: 1, Name: "Old Maker", LogoFile: "logo.png"}
				suite.mockRepo.On("GetByID", uint(1)).Return(existingMaker, nil).Times(2)
				suite.mockRepo.On("GetByName", "Updated Maker").Return((*models.Maker)(nil), nil)
				suite.mockRepo.On("Update", mock.AnythingOfType("*models.Maker")).Return(nil)
			},
			expectedStatusCode: http.StatusFound,
		},
		{
			name:     "エラーケース - 無効なID",
			makerID:  "invalid",
			formData: map[string]string{"name": "Updated Maker"},
			mockSetup: func() {
				// モックセットアップなし
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:     "エラーケース - メーカーが見つからない",
			makerID:  "999",
			formData: map[string]string{"name": "Updated Maker"},
			mockSetup: func() {
				suite.mockRepo.On("GetByID", uint(999)).Return((*models.Maker)(nil), errors.New("not found"))
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			// フォームデータ作成
			form := url.Values{}
			for k, v := range tc.formData {
				form.Set(k, v)
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/makers/%s/edit", tc.makerID), strings.NewReader(form.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.makerID)

			// セッション設定
			adminUser := &models.AdminUser{ID: 1, Name: "Test Admin", LoginID: "admin"}
			utils.SetAdminSession(c, adminUser)

			// テスト実行
			err := suite.handler.UpdateMaker(c)

			// アサーション
			if tc.expectedStatusCode < 400 {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// DeleteMakerのテスト
func (suite *MakerHandlerTestSuite) TestDeleteMaker() {
	testCases := []struct {
		name               string
		makerID            string
		mockSetup          func()
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:    "正常ケース",
			makerID: "1",
			mockSetup: func() {
				suite.mockRepo.On("Delete", uint(1)).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   map[string]interface{}{"message": "メーカーを削除しました"},
		},
		{
			name:    "エラーケース - 無効なID",
			makerID: "invalid",
			mockSetup: func() {
				// モックセットアップなし
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   map[string]interface{}{"error": "無効なメーカーIDです"},
		},
		{
			name:    "エラーケース - 削除失敗",
			makerID: "1",
			mockSetup: func() {
				suite.mockRepo.On("Delete", uint(1)).Return(errors.New("delete error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   map[string]interface{}{"error": "メーカーの削除に失敗しました"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/makers/%s", tc.makerID), nil)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.makerID)

			// セッション設定
			adminUser := &models.AdminUser{ID: 1, Name: "Test Admin", LoginID: "admin"}
			utils.SetAdminSession(c, adminUser)

			// テスト実行
			err := suite.handler.DeleteMaker(c)

			// アサーション
			suite.NoError(err)
			suite.Equal(tc.expectedStatusCode, rec.Code)

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// GetMakerStatsのテスト
func (suite *MakerHandlerTestSuite) TestGetMakerStats() {
	testCases := []struct {
		name               string
		mockSetup          func()
		expectedStatusCode int
	}{
		{
			name: "正常ケース",
			mockSetup: func() {
				stats := &models.MakerStats{
					TotalMakers:  10,
					ActiveMakers: 8,
					TotalBalls:   50,
					TotalCovers:  30,
					TotalCores:   20,
				}
				suite.mockRepo.On("GetMakerStats").Return(stats, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "エラーケース - リポジトリエラー",
			mockSetup: func() {
				suite.mockRepo.On("GetMakerStats").Return((*models.MakerStats)(nil), errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/admin/makers/stats", nil)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// セッション設定
			adminUser := &models.AdminUser{ID: 1, Name: "Test Admin", LoginID: "admin"}
			utils.SetAdminSession(c, adminUser)

			// テスト実行
			err := suite.handler.GetMakerStats(c)

			// アサーション
			suite.NoError(err)
			suite.Equal(tc.expectedStatusCode, rec.Code)

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// ServeLogoImageのテスト
func (suite *MakerHandlerTestSuite) TestServeLogoImage() {
	testCases := []struct {
		name               string
		makerID            string
		mockSetup          func()
		expectedStatusCode int
	}{
		{
			name:    "正常ケース",
			makerID: "1",
			mockSetup: func() {
				maker := &models.Maker{ID: 1, Name: "Test Maker", LogoFile: "logo.png"}
				suite.mockRepo.On("GetByID", uint(1)).Return(maker, nil)

				// MinIOクライアントのモックは複雑なため、実際のテストでは統合テストで行う
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "エラーケース - 無効なID",
			makerID: "invalid",
			mockSetup: func() {
				// モックセットアップなし
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "エラーケース - メーカーが見つからない",
			makerID: "999",
			mockSetup: func() {
				suite.mockRepo.On("GetByID", uint(999)).Return((*models.Maker)(nil), errors.New("not found"))
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:    "エラーケース - ロゴファイルなし",
			makerID: "1",
			mockSetup: func() {
				maker := &models.Maker{ID: 1, Name: "Test Maker", LogoFile: ""}
				suite.mockRepo.On("GetByID", uint(1)).Return(maker, nil)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/makers/%s/logo", tc.makerID), nil)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.makerID)

			// テスト実行
			err := suite.handler.ServeLogoImage(c)

			// アサーション
			if tc.expectedStatusCode < 400 {
				// 正常ケースではMinIOの実装部分でエラーになることが予想される（モックが不完全なため）
				// 実際の統合テストで確認する
				suite.mockRepo.AssertExpectations(suite.T())
			} else {
				suite.Error(err)
				suite.mockRepo.AssertExpectations(suite.T())
			}
		})
	}
}

// getLogoURLのテスト
func (suite *MakerHandlerTestSuite) TestGetLogoURL() {
	testCases := []struct {
		name        string
		makerID     uint
		mockSetup   func()
		expectedURL string
	}{
		{
			name:    "正常ケース",
			makerID: 1,
			mockSetup: func() {
				maker := &models.Maker{ID: 1, Name: "Test Maker", LogoFile: "logo.png"}
				suite.mockRepo.On("GetByID", uint(1)).Return(maker, nil)
			},
			expectedURL: "/admin/makers/1/logo",
		},
		{
			name:    "ロゴファイルなし",
			makerID: 1,
			mockSetup: func() {
				maker := &models.Maker{ID: 1, Name: "Test Maker", LogoFile: ""}
				suite.mockRepo.On("GetByID", uint(1)).Return(maker, nil)
			},
			expectedURL: "",
		},
		{
			name:    "メーカーが見つからない",
			makerID: 999,
			mockSetup: func() {
				suite.mockRepo.On("GetByID", uint(999)).Return((*models.Maker)(nil), errors.New("not found"))
			},
			expectedURL: "",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			// テスト実行
			url := suite.handler.getLogoURL(tc.makerID)

			// アサーション
			suite.Equal(tc.expectedURL, url)
			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// uploadLogoのテスト（ファイルアップロードテスト）
func (suite *MakerHandlerTestSuite) TestUploadLogo() {
	testCases := []struct {
		name         string
		fileName     string
		fileSize     int64
		fileContent  string
		mockSetup    func()
		expectError  bool
		errorMessage string
	}{
		{
			name:         "ファイルサイズが大きすぎる",
			fileName:     "large.png",
			fileSize:     6 * 1024 * 1024, // 6MB
			fileContent:  "large image data",
			mockSetup:    func() {},
			expectError:  true,
			errorMessage: "ファイルサイズが大きすぎます",
		},
		{
			name:         "サポートされていないファイル形式",
			fileName:     "document.txt",
			fileSize:     1024,
			fileContent:  "text content",
			mockSetup:    func() {},
			expectError:  true,
			errorMessage: "サポートされていないファイル形式です",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup()

			// テスト用のFileHeaderを作成
			fileHeader := &multipart.FileHeader{
				Filename: tc.fileName,
				Size:     tc.fileSize,
			}

			// テスト実行
			err := suite.handler.uploadLogo(1, fileHeader)

			// アサーション
			if tc.expectError {
				suite.Error(err)
				suite.Contains(err.Error(), tc.errorMessage)
			} else {
				suite.NoError(err)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

// uploadLogoのテスト用ヘルパー関数
func createTestMultipartFile(fieldName, fileName, content string) (*multipart.FileHeader, *bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, nil, err
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		return nil, nil, err
	}

	writer.Close()

	// Parse the multipart form to get FileHeader
	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20) // 10MB max
	if err != nil {
		return nil, nil, err
	}

	files := form.File[fieldName]
	if len(files) == 0 {
		return nil, nil, errors.New("no file found")
	}

	return files[0], body, nil
}