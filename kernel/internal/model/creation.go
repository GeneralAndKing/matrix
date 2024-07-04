package model

import (
	"gorm.io/gorm"
	"kernel/internal/model/dto"
	"kernel/internal/model/enum"
	"time"
)

type Creation struct {
	gorm.Model
	Type        enum.Creation        `gorm:"type:integer;not null"`
	Title       string               `gorm:"type:varchar(255);not null"`
	Description string               `gorm:"type:text;not null"`
	Paths       GenericArray[string] `gorm:"type:text"`
}

func (w Creation) Output() dto.CreationOutput {
	return dto.CreationOutput{
		Model:       w.Model,
		Type:        w.Type,
		Title:       w.Title,
		Description: w.Description,
		Paths:       w.Paths,
	}
}

type DouyinCreation struct {
	gorm.Model

	DouyinUser        DouyinUser
	DouyinUserId      uint
	Creation          Creation
	CreationID        uint
	VideoId           string                     `gorm:"type:varchar(255)"`
	Title             string                     `gorm:"type:varchar(255);not null"`
	Description       string                     `gorm:"type:text;not null"`
	VideoCoverPath    string                     `gorm:"type:varchar(255)"`
	Location          string                     `gorm:"type:varchar(255)"`
	Flashmob          string                     `gorm:"type:varchar(255)"`
	CollectionName    string                     `gorm:"type:varchar(255)"`
	CollectionNum     int                        `gorm:"type:int(11)"`
	AssociatedHotspot string                     `gorm:"type:varchar(255)"`
	SyncToToutiao     bool                       `gorm:"type:tinyint(1)"`
	AllowedToSave     bool                       `gorm:"type:tinyint(1);not null"`
	WhoCanWatch       enum.WhoCanWatch           `gorm:"type:integer;not null"`
	ReleaseTime       time.Time                  `gorm:"type:datetime;not null"`
	Status            enum.PublishCreationStatus `gorm:"type:integer;not null"`
	Message           string                     `gorm:"type:text;not null"`
}

func (c DouyinCreation) Output() dto.DouyinCreationOutput {
	return dto.DouyinCreationOutput{
		Model:             c.Model,
		DouyinUserId:      c.DouyinUserId,
		CreationID:        c.CreationID,
		Title:             c.Title,
		Description:       c.Description,
		VideoCoverPath:    c.VideoCoverPath,
		Location:          c.Location,
		Paster:            c.Flashmob,
		CollectionName:    c.CollectionName,
		CollectionNum:     c.CollectionNum,
		AssociatedHotspot: c.AssociatedHotspot,
		SyncToToutiao:     c.SyncToToutiao,
		AllowedToSave:     c.AllowedToSave,
		WhoCanWatch:       c.WhoCanWatch,
		ReleaseTime:       c.ReleaseTime,
		Status:            c.Status,
	}
}
