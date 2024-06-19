package dto

import (
	"gorm.io/gorm"
	"kernel/internal/model/enum"
)

type LabelOutput string

type WorkOutput struct {
	gorm.Model
	Type        enum.Work `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Paths       []string  `json:"paths"`
}

type DouyinUserOutput struct {
	gorm.Model
	Name        string        `json:"name"`
	Description string        `json:"description"`
	DouyinId    string        `json:"douyinId"`
	Avatar      string        `json:"avatar"`
	Labels      []LabelOutput `json:"labels"`
	Expired     bool          `json:"expired"`
}
