package model

import (
	"gorm.io/gorm"
	"kernel/internal/model/dto"
	"kernel/internal/model/enum"
	"time"
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
	DouyinUser        DouyinUser
	DouyinUserId      uint
	Work              Work
	WorkID            uint
	Title             string                 `gorm:"type:varchar(255);not null"`
	Description       string                 `gorm:"type:text;not null"`
	VideoCoverPath    string                 `gorm:"type:varchar(255)"`
	Location          string                 `gorm:"type:varchar(255)"`
	Paster            string                 `gorm:"type:varchar(255)"`
	CollectionName    string                 `gorm:"type:varchar(255)"`
	CollectionNum     int                    `gorm:"type:int(11)"`
	AssociatedHotspot string                 `gorm:"type:varchar(255)"`
	SyncToToutiao     bool                   `gorm:"type:tinyint(1)"`
	AllowedToSave     bool                   `gorm:"type:tinyint(1);not null"`
	WhoCanWatch       uint                   `gorm:"type:integer;not null"`
	ReleaseTime       time.Time              `gorm:"type:datetime;not null"`
	Status            enum.PublishWorkStatus `gorm:"type:integer;not null"`
	Message           string                 `gorm:"type:text;not null"`
}
