package dto

import (
	"gorm.io/gorm"
	"kernel/internal/model/enum"
)

type WorkOutput struct {
	gorm.Model
	Type        enum.Work `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Paths       []string  `json:"paths"`
}
