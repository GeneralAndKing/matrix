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

type WorkPageInput struct {
	Types []enum.Work `form:"type" binding:"dive,min=1,max=2"`
	Size  int         `form:"size,default=10"`
	Page  int         `form:"page,default=0"`
}
