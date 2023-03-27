package models

import "gorm.io/gorm"

type TaskUser struct {
	gorm.Model
	TaskID      uint   `gorm:"not null;foreignKey" json:"task_id"`
	UserID      uint   `gorm:"not null;foreignKey" json:"user_id"`
	Description string `json:"description"`
	Status      uint8  `gorm:"not null;size:1;default:0;comment:0 => on progress, 1 => completed" json:"status"`
}

type TaskUserRequest struct {
	TaskID      uint   `json:"task_id"`
	UserID      uint   `json:"user_id"`
	Description string `json:"description"`
	Status      uint8  `json:"status"`
}

type TaskUserResponse struct {
	User        UserResponse `json:"user"`
	Description string       `json:"description"`
	Status      uint8        `json:"status"`
}

type UserTaskResponse struct {
	Task        TaskResponse `json:"task"`
	Description string       `json:"description"`
	Status      uint8        `json:"status"`
}
