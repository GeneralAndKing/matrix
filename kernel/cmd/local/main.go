package main

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kernel/internal/api"
	"kernel/internal/config"
	"kernel/internal/database"
	"kernel/internal/model"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
}

func main() {
	initLog()
	initDatabase()
	initHttpServer()
	_ = zap.L().Sync()
}
