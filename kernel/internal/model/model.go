package model

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(new(Work), new(DouyinUser), new(Label))
}
