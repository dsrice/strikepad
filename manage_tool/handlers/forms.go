package handlers

import (
	"errors"
	"strconv"
	"strings"
)

// CoverForm はカバー作成・編集フォーム用の構造体
type CoverForm struct {
	Name         string `form:"name" validate:"required"`
	MakerID      uint   `form:"maker_id" validate:"required"`
	MaterialType int    `form:"material_type" validate:"required,min=1,max=5"`
	Rank         int    `form:"rank" validate:"required,min=1,max=99"`
}

// Validate はカスタムバリデーションを行う
func (f *CoverForm) Validate() error {
	// 名前のトリムと検証
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return errors.New("カバー名は必須です")
	}

	// メーカーIDの検証
	if f.MakerID == 0 {
		return errors.New("メーカーの選択は必須です")
	}

	// 材質種別の検証
	if f.MaterialType < 1 || f.MaterialType > 5 {
		return errors.New("無効な材質種別です（1-5の範囲で選択してください）")
	}

	// カバーの強さの検証
	if f.Rank < 1 || f.Rank > 99 {
		return errors.New("無効なカバーの強さです（1-99の範囲で入力してください）")
	}

	return nil
}

// SearchForm は検索フォーム用の構造体
type SearchForm struct {
	Search  string `query:"search"`
	MakerID string `query:"maker_id"`
	Page    int    `query:"page"`
}

// ParseMakerID はMakerIDを数値に変換する
func (f *SearchForm) ParseMakerID() (uint, error) {
	if f.MakerID == "" {
		return 0, nil
	}

	id, err := strconv.ParseUint(f.MakerID, 10, 32)
	if err != nil {
		return 0, errors.New("無効なメーカーIDです")
	}

	return uint(id), nil
}

// GetValidPage は有効なページ番号を返す
func (f *SearchForm) GetValidPage() int {
	if f.Page <= 0 {
		return 1
	}
	return f.Page
}
