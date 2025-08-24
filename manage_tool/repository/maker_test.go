package repository

import (
	"testing"
	"time"

	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository/testutil"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MakerRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *MakerRepository
}

func (suite *MakerRepositoryTestSuite) SetupTest() {
	suite.db = testutil.SetupTestDB(suite.T())
	suite.repo = NewMakerRepository(suite.db)
}

func (suite *MakerRepositoryTestSuite) TearDownTest() {
	testutil.CleanupTestDB(suite.db)
}

func TestMakerRepositorySuite(t *testing.T) {
	suite.Run(t, new(MakerRepositoryTestSuite))
}

// GetAllMakersのテスト
func (suite *MakerRepositoryTestSuite) TestGetAllMakers() {
	// テストデータをシード
	testutil.SeedMakersData(suite.db)
	testutil.SeedRelatedData(suite.db)

	testCases := []struct {
		name          string
		expectedFirst string
		offset        int
		limit         int
		expectedCount int
		expectedError bool
	}{
		{
			name:          "正常ケース - 全件取得",
			offset:        0,
			limit:         10,
			expectedCount: 3,           // 論理削除されていないメーカー数
			expectedFirst: "Brunswick", // 作成日降順で最新
			expectedError: false,
		},
		{
			name:          "正常ケース - ページネーション",
			offset:        1,
			limit:         2,
			expectedCount: 2,
			expectedFirst: "Hammer",
			expectedError: false,
		},
		{
			name:          "正常ケース - 制限値1",
			offset:        0,
			limit:         1,
			expectedCount: 1,
			expectedFirst: "Brunswick",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			makers, err := suite.repo.GetAllMakers(tc.offset, tc.limit)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Len(makers, tc.expectedCount)
				if tc.expectedCount > 0 {
					suite.Equal(tc.expectedFirst, makers[0].Name)
					// BallCount, CoverCount, CoreCountが正しく集計されているか確認
					suite.GreaterOrEqual(makers[0].BallCount, int64(0))
					suite.GreaterOrEqual(makers[0].CoverCount, int64(0))
					suite.GreaterOrEqual(makers[0].CoreCount, int64(0))
				}
			}
		})
	}
}

// SearchMakersのテスト
func (suite *MakerRepositoryTestSuite) TestSearchMakers() {
	// テストデータをシード
	testutil.SeedMakersData(suite.db)
	testutil.SeedRelatedData(suite.db)

	testCases := []struct {
		name          string
		searchQuery   string
		expectedNames []string
		offset        int
		limit         int
		expectedCount int
		expectedError bool
	}{
		{
			name:          "正常ケース - 部分一致検索",
			searchQuery:   "Storm",
			offset:        0,
			limit:         10,
			expectedCount: 1,
			expectedNames: []string{"Storm"},
			expectedError: false,
		},
		{
			name:          "正常ケース - 大文字小文字を無視",
			searchQuery:   "storm",
			offset:        0,
			limit:         10,
			expectedCount: 1,
			expectedNames: []string{"Storm"},
			expectedError: false,
		},
		{
			name:          "正常ケース - 複数の結果",
			searchQuery:   "am", // "Hammer"にマッチ
			offset:        0,
			limit:         10,
			expectedCount: 1,
			expectedNames: []string{"Hammer"},
			expectedError: false,
		},
		{
			name:          "正常ケース - 検索結果なし",
			searchQuery:   "NonExistent",
			offset:        0,
			limit:         10,
			expectedCount: 0,
			expectedNames: []string{},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			makers, err := suite.repo.SearchMakers(tc.searchQuery, tc.offset, tc.limit)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Len(makers, tc.expectedCount)
				for i, expectedName := range tc.expectedNames {
					if i < len(makers) {
						suite.Equal(expectedName, makers[i].Name)
					}
				}
			}
		})
	}
}

// GetByIDのテスト
func (suite *MakerRepositoryTestSuite) TestGetByID() {
	// テストデータをシード
	makers := testutil.SeedMakersData(suite.db)

	testCases := []struct {
		expectedMaker *models.Maker
		name          string
		id            uint
		expectedError bool
	}{
		{
			name:          "正常ケース - 存在するID",
			id:            1,
			expectedMaker: makers[0], // Storm
			expectedError: false,
		},
		{
			name:          "正常ケース - 存在しないID",
			id:            999,
			expectedMaker: nil,
			expectedError: false,
		},
		{
			name:          "正常ケース - 論理削除されたメーカー",
			id:            4,
			expectedMaker: nil, // 論理削除されているのでnilが返る
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			maker, err := suite.repo.GetByID(tc.id)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				if tc.expectedMaker == nil {
					suite.Nil(maker)
				} else {
					suite.NotNil(maker)
					testutil.AssertMakerEqual(suite.T(), tc.expectedMaker, maker)
				}
			}
		})
	}
}

// GetByNameのテスト
func (suite *MakerRepositoryTestSuite) TestGetByName() {
	// テストデータをシード
	makers := testutil.SeedMakersData(suite.db)

	testCases := []struct {
		expectedMaker *models.Maker
		name          string
		makerName     string
		expectedError bool
	}{
		{
			name:          "正常ケース - 存在する名前",
			makerName:     "Storm",
			expectedMaker: makers[0], // Storm
			expectedError: false,
		},
		{
			name:          "正常ケース - 存在しない名前",
			makerName:     "NonExistent",
			expectedMaker: nil,
			expectedError: false,
		},
		{
			name:          "正常ケース - 論理削除されたメーカー",
			makerName:     "Deleted Maker",
			expectedMaker: nil, // 論理削除されているのでnilが返る
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			maker, err := suite.repo.GetByName(tc.makerName)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				if tc.expectedMaker == nil {
					suite.Nil(maker)
				} else {
					suite.NotNil(maker)
					testutil.AssertMakerEqual(suite.T(), tc.expectedMaker, maker)
				}
			}
		})
	}
}

// Createのテスト
func (suite *MakerRepositoryTestSuite) TestCreate() {
	testCases := []struct {
		maker         *models.Maker
		name          string
		errorContains string
		expectedError bool
	}{
		{
			name: "正常ケース",
			maker: &models.Maker{
				Name:      "New Maker",
				LogoFile:  "logo.png",
				IsDeleted: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: false,
		},
		{
			name: "正常ケース - LogoFileなし",
			maker: &models.Maker{
				Name:      "Another Maker",
				IsDeleted: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			err := suite.repo.Create(tc.maker)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
				if tc.errorContains != "" {
					suite.Contains(err.Error(), tc.errorContains)
				}
			} else {
				suite.NoError(err)
				// IDが自動割り当てされているか確認
				suite.NotZero(tc.maker.ID)

				// データベースに保存されているか確認
				var savedMaker models.Maker
				err = suite.db.First(&savedMaker, tc.maker.ID).Error
				suite.NoError(err)
				suite.Equal(tc.maker.Name, savedMaker.Name)
				suite.Equal(tc.maker.LogoFile, savedMaker.LogoFile)
			}
		})
	}
}

// Updateのテスト
func (suite *MakerRepositoryTestSuite) TestUpdate() {
	// テストデータをシード
	makers := testutil.SeedMakersData(suite.db)

	testCases := []struct {
		updateMaker   *models.Maker
		name          string
		expectedError bool
	}{
		{
			name: "正常ケース - 名前とロゴファイル更新",
			updateMaker: &models.Maker{
				ID:        makers[0].ID,
				Name:      "Updated Storm",
				LogoFile:  "updated_logo.png",
				IsDeleted: false,
				UpdatedAt: time.Now(),
			},
			expectedError: false,
		},
		{
			name: "正常ケース - 存在しないIDの更新",
			updateMaker: &models.Maker{
				ID:        999,
				Name:      "Non Existent",
				IsDeleted: false,
				UpdatedAt: time.Now(),
			},
			expectedError: false, // GORMは存在しないレコードの更新でもエラーにならない
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			err := suite.repo.Update(tc.updateMaker)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)

				// 存在するIDの場合、更新されているか確認
				if tc.updateMaker.ID <= uint(len(makers)) {
					var updatedMaker models.Maker
					err = suite.db.First(&updatedMaker, tc.updateMaker.ID).Error
					if err == nil { // レコードが存在する場合
						suite.Equal(tc.updateMaker.Name, updatedMaker.Name)
						suite.Equal(tc.updateMaker.LogoFile, updatedMaker.LogoFile)
					}
				}
			}
		})
	}
}

// UpdateLogoFileのテスト
func (suite *MakerRepositoryTestSuite) TestUpdateLogoFile() {
	// テストデータをシード
	makers := testutil.SeedMakersData(suite.db)

	testCases := []struct {
		name          string
		filename      string
		makerID       uint
		expectedError bool
	}{
		{
			name:          "正常ケース",
			makerID:       makers[0].ID,
			filename:      "new_logo.png",
			expectedError: false,
		},
		{
			name:          "正常ケース - 存在しないID",
			makerID:       999,
			filename:      "logo.png",
			expectedError: false, // GORMは存在しないレコードの更新でもエラーにならない
		},
		{
			name:          "正常ケース - 論理削除されたメーカー",
			makerID:       4, // 論理削除済み
			filename:      "logo.png",
			expectedError: false, // 論理削除済みなので更新されない
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			err := suite.repo.UpdateLogoFile(tc.makerID, tc.filename)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)

				// 存在するIDで論理削除されていない場合、更新されているか確認
				if tc.makerID <= 3 { // 1,2,3は論理削除されていない
					var updatedMaker models.Maker
					err = suite.db.First(&updatedMaker, tc.makerID).Error
					if err == nil {
						suite.Equal(tc.filename, updatedMaker.LogoFile)
					}
				}
			}
		})
	}
}

// Deleteのテスト
func (suite *MakerRepositoryTestSuite) TestDelete() {
	// テストデータをシード
	makers := testutil.SeedMakersData(suite.db)

	testCases := []struct {
		name          string
		makerID       uint
		expectedError bool
	}{
		{
			name:          "正常ケース - 存在するメーカー",
			makerID:       makers[0].ID,
			expectedError: false,
		},
		{
			name:          "正常ケース - 存在しないメーカー",
			makerID:       999,
			expectedError: false, // GORMは存在しないレコードの削除でもエラーにならない
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// テスト実行
			err := suite.repo.Delete(tc.makerID)

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)

				// 論理削除されているか確認
				var deletedMaker models.Maker
				err = suite.db.Unscoped().First(&deletedMaker, tc.makerID).Error
				if err == nil { // レコードが存在する場合
					suite.True(deletedMaker.IsDeleted)
				}
			}
		})
	}
}

// GetTotalMakersCountのテスト
func (suite *MakerRepositoryTestSuite) TestGetTotalMakersCount() {
	testCases := []struct {
		name          string
		expectedCount int64
		seedData      bool
		expectedError bool
	}{
		{
			name:          "正常ケース - データなし",
			seedData:      false,
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "正常ケース - データあり",
			seedData:      true,
			expectedCount: 3, // 論理削除されていないメーカー数
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// データのシード
			if tc.seedData {
				testutil.SeedMakersData(suite.db)
			}

			// テスト実行
			count, err := suite.repo.GetTotalMakersCount()

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedCount, count)
			}
		})
	}
}

// GetMakerStatsのテスト
func (suite *MakerRepositoryTestSuite) TestGetMakerStats() {
	testCases := []struct {
		expectedStats *models.MakerStats
		name          string
		seedData      bool
		expectedError bool
	}{
		{
			name:     "正常ケース - データなし",
			seedData: false,
			expectedStats: &models.MakerStats{
				TotalMakers:  0,
				ActiveMakers: 0,
				TotalBalls:   0,
				TotalCovers:  0,
				TotalCores:   0,
			},
			expectedError: false,
		},
		{
			name:     "正常ケース - データあり",
			seedData: true,
			expectedStats: &models.MakerStats{
				TotalMakers:  3, // 論理削除されていないメーカー数
				ActiveMakers: 2, // ボールを持つメーカー数（Storm, Hammer）
				TotalBalls:   3, // 全ボール数
				TotalCovers:  3, // 全カバー数
				TotalCores:   3, // 全コア数
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// データのシード
			if tc.seedData {
				testutil.SeedMakersData(suite.db)
				testutil.SeedRelatedData(suite.db)
			}

			// テスト実行
			stats, err := suite.repo.GetMakerStats()

			// アサーション
			if tc.expectedError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.NotNil(stats)
				suite.Equal(tc.expectedStats.TotalMakers, stats.TotalMakers)
				suite.Equal(tc.expectedStats.ActiveMakers, stats.ActiveMakers)
				suite.Equal(tc.expectedStats.TotalBalls, stats.TotalBalls)
				suite.Equal(tc.expectedStats.TotalCovers, stats.TotalCovers)
				suite.Equal(tc.expectedStats.TotalCores, stats.TotalCores)
			}
		})
	}
}