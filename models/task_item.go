package models

import "gorm.io/gorm"

type TaskItem struct {
	gorm.Model
	TaskID uint   `gorm:"not null;foreignKey" json:"task_id"`
	Name   string `gorm:"not null" json:"name"`
}

type TaskItemResponse struct {
	Name string `json:"name"`
}
