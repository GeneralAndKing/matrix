package dto

import (
	"kernel/internal/model/enum"
	"time"
)

type AddCreationInput struct {
	Type        enum.Creation `json:"type" binding:"min=1,max=2"`
	Title       string        `json:"title" binding:"required"`
	Description string        `json:"description" binding:"required"`
	Paths       []string      `json:"paths" binding:"required"`
}

type AddDouyinUserInput []string

type PublishCreationInput struct {
	DouyinUserInputs []PublishDouyinCreationInput `json:"douyin" binding:"required"`
}

type PublishDouyinCreationInput struct {
	ID                string           `json:"id" binding:"required"`
	Title             string           `json:"title" binding:"required"`
	Description       string           `json:"description" binding:"required"`
	VideoCoverPath    string           `json:"videoCoverPath"`
	Location          string           `json:"location"`
	Paster            string           `json:"paster"`
	CollectionName    string           `json:"collectionName"`
	CollectionNum     int              `json:"collectionNum"`
	AssociatedHotspot string           `json:"associatedHotspot"`
	SyncToToutiao     bool             `json:"syncToToutiao" binding:"required"`
	AllowedToSave     bool             `json:"allowedToSave" binding:"required"`
	WhoCanWatch       enum.WhoCanWatch `json:"whoCanWatch" binding:"required"`
	ReleaseTime       time.Time        `json:"releaseTime" binding:"required"`
}
