package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kernel/internal/api"
	"kernel/internal/config"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/ws"
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

func initDatabase(ctx context.Context) {
	dataSource := config.GetDataSource()
	database.RegisterSqlite3(dataSource.Sqlite3)
	if err := database.Sqlite3Manager(ctx, func(db *gorm.DB) error {
		return model.Migrate(db)
	}); err != nil {
		panic(err)
	}

}

func notifyExitSignal(ctx context.Context) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		zap.L().Info("ctx exit is detected", zap.Error(ctx.Err()))
	case s := <-quit:
		zap.L().Info("exit signal detected", zap.String("signal", s.String()))
	}
}

func initHttpServer(ctx context.Context) {
	application := config.GetApplication()
	ws.Init()
	handler := api.API(ctx, application.Debug)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		zap.L().Info("start ws")
		if err := ws.Open(); err != nil {
			zap.L().Fatal("ws error", zap.Error(err))
		}
	}()
	defer ws.Close()
	go func() {
		zap.L().Info("start server")
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("server error", zap.Error(err))
		}
	}()
	notifyExitSignal(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("server shutdown error", zap.Error(err))
	}
}

func main() {
	ctx := context.Background()
	initLog()
	initDatabase(ctx)
	initHttpServer(ctx)
	_ = zap.L().Sync()
}
