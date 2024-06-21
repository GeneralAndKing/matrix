package dto

import (
	"kernel/internal/model/enum"
	"time"
)

type AddWorkInput struct {
	Type        enum.Work `json:"type" binding:"min=1,max=2"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Paths       []string  `json:"paths" binding:"required"`
}

type AddDouyinUserInput []string

type PublishWorkInput struct {
	DouyinUserInputs []PublishDouyinWorkInput `json:"douyin" binding:"required"`
}

type PublishDouyinWorkInput struct {
	ID                string    `json:"id" binding:"required"`
	Title             string    `json:"title" binding:"required"`
	Description       string    `json:"description" binding:"required"`
	VideoCoverPath    string    `json:"video_cover_path"`
	Location          string    `json:"location"`
	Paster            string    `json:"paster"`
	CollectionName    string    `json:"collection_name"`
	CollectionNum     int       `json:"collection_num"`
	AssociatedHotspot string    `json:"associated_hotspot"`
	SyncToToutiao     bool      `json:"sync_to_toutiao" binding:"required"`
	AllowedToSave     bool      `json:"allowed_to_save" binding:"required"`
	WhoCanWatch       uint      `json:"who_can_watch" binding:"required"`
	ReleaseTime       time.Time `json:"release_time" binding:"required"`
}
