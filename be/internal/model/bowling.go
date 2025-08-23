package model

import (
	"time"

	"gorm.io/gorm"
)

// Maker はメーカーモデル
type Maker struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Covers    []Cover        `gorm:"foreignKey:MakerID" json:"covers,omitempty"`
	Cores     []Core         `gorm:"foreignKey:MakerID" json:"cores,omitempty"`
	Balls     []Ball         `gorm:"foreignKey:MakerID" json:"balls,omitempty"`
	ID        uint           `gorm:"primarykey" json:"id"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Maker) TableName() string {
	return "makers"
}

// Cover はカバーモデル
type Cover struct {
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	Maker        Maker          `gorm:"foreignKey:MakerID" json:"maker,omitempty"`
	Balls        []Ball         `gorm:"foreignKey:CoverID" json:"balls,omitempty"`
	ID           uint           `gorm:"primarykey" json:"id"`
	MaterialType int            `gorm:"not null" json:"material_type"`
	MakerID      uint           `gorm:"not null" json:"maker_id"`
	Rank         int            `gorm:"not null" json:"rank"`
	IsDeleted    bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Cover) TableName() string {
	return "covers"
}

// Core はコアモデル
type Core struct {
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	InitDiff     *float32       `json:"init_diff,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	Maker        Maker          `gorm:"foreignKey:MakerID" json:"maker,omitempty"`
	Balls        []Ball         `gorm:"foreignKey:CoreID" json:"balls,omitempty"`
	ID           uint           `gorm:"primarykey" json:"id"`
	MakerID      uint           `gorm:"not null" json:"maker_id"`
	DeltaRG      float32        `gorm:"not null" json:"delta_rg"`
	RG           float32        `gorm:"not null" json:"rg"`
	IsDeleted    bool           `gorm:"default:false" json:"is_deleted"`
	SymmetryFlag bool           `gorm:"not null" json:"symmetry_flag"`
}

// TableName はテーブル名を指定
func (Core) TableName() string {
	return "cores"
}

// Ball はボールモデル
type Ball struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	URL       string         `gorm:"size:200;not null" json:"url"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Core      Core           `gorm:"foreignKey:CoreID" json:"core,omitempty"`
	Cover     Cover          `gorm:"foreignKey:CoverID" json:"cover,omitempty"`
	Maker     Maker          `gorm:"foreignKey:MakerID" json:"maker,omitempty"`
	CoreID    uint           `gorm:"not null" json:"core_id"`
	CoverID   uint           `gorm:"not null" json:"cover_id"`
	MakerID   uint           `gorm:"not null" json:"maker_id"`
	ID        uint           `gorm:"primarykey" json:"id"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Ball) TableName() string {
	return "balls"
}

// MakerListItem はメーカー一覧表示用の軽量構造体
type MakerListItem struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `json:"name"`
	ID         uint      `json:"id"`
	BallCount  int64     `json:"ball_count"`
	CoverCount int64     `json:"cover_count"`
	CoreCount  int64     `json:"core_count"`
}

// MakerStats はメーカー統計用構造体
type MakerStats struct {
	TotalMakers  int64 `json:"total_makers"`
	ActiveMakers int64 `json:"active_makers"`
	TotalBalls   int64 `json:"total_balls"`
	TotalCovers  int64 `json:"total_covers"`
	TotalCores   int64 `json:"total_cores"`
}