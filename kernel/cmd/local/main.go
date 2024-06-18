package main

import (
	"context"
	"github.com/fvbock/endless"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kernel/internal/api"
	"kernel/internal/config"
	"kernel/internal/database"
	"kernel/internal/model"
	"log"
)

type zapWriter struct {
	logger *zap.Logger
}

func (w zapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}

func initLog() {
	application := config.GetApplication()
	if application.Debug {

	}
	production, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(production)
	log.SetOutput(zapWriter{logger: production})
}

func initDatabase() {
	dataSource := config.GetDataSource()
	database.RegisterSqlite3(dataSource.Sqlite3)
	if err := database.Sqlite3Manager(context.Background(), func(db *gorm.DB) error {
		return model.Migrate(db)
	}); err != nil {
		panic(err)
	}

}

func initHttpServer() {
	application := config.GetApplication()
	handler := api.API(application.Debug)
	if err := endless.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}

func main() {
	initLog()
	initDatabase()
	initHttpServer()
	_ = zap.L().Sync()
}
