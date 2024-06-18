package api

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kernel/internal/service"
	"net/http"
	"time"
)

func API(debug bool) http.Handler {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(zap.L(), true))

	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	workGroup := engine.Group("/work")
	workGroup.GET("/", service.GetWorkPage)

	return engine
}
