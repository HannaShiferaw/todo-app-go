package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	Priority    string     `json:"priority" gorm:"default:medium"` // low, medium, high
	DueDate     *time.Time `json:"due_date"`
	Category    string     `json:"category" gorm:"default:personal"` // work, personal, shopping, etc.
}
