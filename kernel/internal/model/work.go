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

func (w Work) ToOutput() dto.WorkOutput {
	return dto.WorkOutput{
		Model:       w.Model,
		Type:        w.Type,
		Title:       w.Title,
		Description: w.Description,
		Paths:       w.Paths,
	}
}

type Topic struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null"`
}
