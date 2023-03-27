package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name         string `gorm:"not null" json:"name"`
	Description  string `json:"description"`
	Score        uint8  `gorm:"not null;size:1" json:"score"`
	Duration     uint8  `gorm:"not null;size:2;default:1;comment:workdays" json:"duration"`
	StartedAt    string `json:"started_at"`
	FinishedAt   string `json:"finished_at"`
	Status       uint8  `gorm:"not null;size:1;default:0;comment:0 => planned, 1 => on progress, 2 => completed" json:"status"`
	TaskDocument TaskDocument
	TaskItems    []TaskItem
	TaskUsers    []TaskUser `json:"users"`
}

type TaskSimpleResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Score      uint8  `json:"score"`
	Duration   uint8  `json:"duration"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	Status     uint8  `json:"status"`
}

type TaskResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Score        uint8  `json:"score"`
	Duration     uint8  `json:"duration"`
	StartedAt    string `json:"started_at"`
	FinishedAt   string `json:"finished_at"`
	Status       uint8  `json:"status"`
	TaskDocument TaskDocumentResponse
	TaskItems    []TaskItemResponse
	TaskUsers    []TaskUserResponse
}

type TaskCreateRequest struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Score        uint8        `json:"score"`
	Duration     uint8        `json:"duration"`
	TaskDocument TaskDocument `json:"task_document"`
	TaskItems    []TaskItem   `json:"task_items"`
}

type TaskUpdateRequest struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Score        uint8        `json:"score"`
	Duration     uint8        `json:"duration"`
	StartedAt    string       `json:"started_at"`
	Status       uint8        `json:"status"`
	TaskDocument TaskDocument `json:"task_document"`
	TaskItems    []TaskItem   `json:"task_items"`
}

type TaskStartedAtRequest struct {
	StartedAt string `json:"started_at"`
}

type TaskStatusRequest struct {
	Status uint8 `json:"status"`
}
