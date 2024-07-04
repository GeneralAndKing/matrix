package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type GenericArray[T any] []T

func (g *GenericArray[T]) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

func (g GenericArray[T]) Value() (driver.Value, error) {
	return json.Marshal(g)

}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(new(Creation), new(DouyinUser), new(Label), new(DouyinCreation), new(DouyinLive))
}
