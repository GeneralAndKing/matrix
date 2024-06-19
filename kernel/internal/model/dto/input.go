package dto

import (
	"kernel/internal/model/enum"
)

type AddWorkInput struct {
	Type        enum.Work `json:"type" binding:"min=1,max=2"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Paths       []string  `json:"paths" binding:"required"`
}

type AddDouyinUserInput []string
