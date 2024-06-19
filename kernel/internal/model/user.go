package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"kernel/internal/model/dto"
	"kernel/pkg/chromedp_ext"
)

type Cookies []chromedp_ext.Cookie

func (c *Cookies) Scan(value interface{}) error {
	v, _ := value.([]byte)
	var receiver Cookies
	err := json.Unmarshal(v, &receiver)
	if err != nil {
		return err
	}
	*c = receiver
	return nil
}

func (c *Cookies) Value() (driver.Value, error) {
	return json.Marshal(c)

}

type Label struct {
	gorm.Model
	Name string `gorm:"varchar(255)"`
}

func (l Label) Output() dto.LabelOutput {
	return dto.LabelOutput(l.Name)

}

type DouyinUser struct {
	gorm.Model
	Name        string  `gorm:"varchar(255)"`
	Description string  `gorm:"varchar(255)"`
	DouyinId    string  `gorm:"varchar(255)"`
	Avatar      string  `gorm:"varchar(1024)"`
	Labels      []Label `gorm:"many2many:douyin_user_labels"`
	Cookies     Cookies `gorm:"json"`
	Expired     bool    `gorm:"bool"`
}

func (d DouyinUser) Output() dto.DouyinUserOutput {
	var labelOutputs []dto.LabelOutput
	for _, label := range d.Labels {
		labelOutputs = append(labelOutputs, label.Output())
	}
	return dto.DouyinUserOutput{
		Model:       d.Model,
		Name:        d.Name,
		Description: d.Description,
		DouyinId:    d.DouyinId,
		Avatar:      d.Avatar,
		Labels:      labelOutputs,
		Expired:     d.Expired,
	}
}
