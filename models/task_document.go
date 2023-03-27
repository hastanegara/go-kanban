package models

import "gorm.io/gorm"

type TaskDocument struct {
	gorm.Model
	TaskID uint   `gorm:"not null;foreignKey" json:"task_id"`
	Path   string `gorm:"not null" json:"path"`
}

type TaskDocumentResponse struct {
	Path string `json:"path"`
}
