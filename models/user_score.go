package models

import "gorm.io/gorm"

type UserScore struct {
	gorm.Model
	UserID uint  `gorm:"not null;foreignKey" json:"user_id"`
	Score  uint8 `gorm:"not null" json:"score"`
}

type UserScoreResponse struct {
	Score uint8 `json:"score"`
}
