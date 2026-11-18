// Package models holds the Todo type stored in the database.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo is one task. gorm.Model adds ID, CreatedAt, UpdatedAt and
// DeletedAt (soft delete). Priority is low, medium or high;
// Category is work, personal, shopping, health or other.
type Todo struct {
	gorm.Model
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	Priority    string     `json:"priority" gorm:"default:medium"`
	DueDate     *time.Time `json:"due_date"`
	Category    string     `json:"category" gorm:"default:personal"`
}
