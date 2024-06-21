package model

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(new(Creation), new(DouyinUser), new(Label), new(DouyinCreation))
}
