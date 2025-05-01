package models

import (
	"time"

	"gorm.io/gorm"
)

// --- Model Definition ---
type Task struct {
	ID          uint   `gorm:"primarykey"`
	Title       string `gorm:"not null"`
	Description string
	Completed   bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// --- DTOs (Data Transfer Objects) ---
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
