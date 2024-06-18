package database

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var (
	sqliteDB *gorm.DB
)

func RegisterSqlite3(dbPath string) {
	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqliteDB = db
}
func Sqlite3Manager(ctx context.Context, handler func(*gorm.DB) error) error {
	err := handler(sqliteDB.WithContext(ctx))
	return err
}

func Sqlite3Transaction(ctx context.Context, handler func(*gorm.DB) error) error {
	return sqliteDB.WithContext(ctx).Transaction(handler)
}
