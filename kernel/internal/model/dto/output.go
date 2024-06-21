package dto

import (
	"gorm.io/gorm"
	"kernel/internal/model/enum"
	"time"
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

type DouyinWorkOutput struct {
	gorm.Model
	DouyinUserId      uint                   `json:"douyinUserId"`
	WorkID            uint                   `json:"workId"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	VideoCoverPath    string                 `json:"videoCoverPath"`
	Location          string                 `json:"location"`
	Paster            string                 `json:"paster"`
	CollectionName    string                 `json:"collectionName"`
	CollectionNum     int                    `json:"collectionNum"`
	AssociatedHotspot string                 `json:"associatedHotspot"`
	SyncToToutiao     bool                   `json:"syncToToutiao"`
	AllowedToSave     bool                   `json:"allowedToSave"`
	WhoCanWatch       uint                   `json:"whoCanWatch"`
	ReleaseTime       time.Time              `json:"releaseTime"`
	Status            enum.PublishWorkStatus `json:"status"`
}
