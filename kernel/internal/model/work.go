package model

import (
	"gorm.io/gorm"
	"kernel/internal/model/dto"
	"kernel/internal/model/enum"
)

type Work struct {
	gorm.Model
	Type        enum.Work `gorm:"type:tinyint(1)"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	Paths       []string  `gorm:"type:text[]"`
}

func (w Work) Output() dto.WorkOutput {
	return dto.WorkOutput{
		Model:       w.Model,
		Type:        w.Type,
		Title:       w.Title,
		Description: w.Description,
		Paths:       w.Paths,
	}
}

type DouyinWork struct {
	gorm.Model
	DouyinId    string
	Work        Work
	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text;not null"`
	//TODO 细化
}
