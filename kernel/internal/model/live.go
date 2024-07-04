package model

import (
	"gorm.io/gorm"
	"kernel/internal/model/dto"
	"kernel/internal/model/enum"
)

type DouyinLive struct {
	gorm.Model
	LiveId   string                 `gorm:"unique;type:varchar(255);not null"`
	Name     string                 `gorm:"type:varchar(255);not null"`
	DouyinId string                 `gorm:"type:varchar(255);not null"`
	Avatar   string                 `gorm:"type:varchar(1024)"`
	Labels   []Label                `gorm:"many2many:douyin_live_labels"`
	Monitor  enum.LiveMonitorStatus `gorm:"type:integer;not null"`
}

func (d DouyinLive) Output() dto.DouyinLiveOutput {
	var labelOutputs []dto.LabelOutput
	for _, label := range d.Labels {
		labelOutputs = append(labelOutputs, label.Output())
	}
	return dto.DouyinLiveOutput{
		Model:    d.Model,
		LiveId:   d.LiveId,
		Name:     d.Name,
		DouyinId: d.DouyinId,
		Avatar:   d.Avatar,
		Labels:   labelOutputs,
		Monitor:  d.Monitor,
	}
}
