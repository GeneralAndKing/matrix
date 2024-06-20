package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
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
	engine.Use(cors.Default())
	engine.Use(func(ctx *gin.Context) {
		ctx.Next()
		if ctx.Errors != nil {
			ctx.JSON(ctx.Writer.Status(), gin.H{
				"error":   http.StatusText(ctx.Writer.Status()),
				"message": ctx.Errors.String(),
			})
		}

	})
	engine.GET("/message")
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	workGroup := engine.Group("/work")
	workGroup.GET("", service.GetAllWork)
	workGroup.POST("", service.AddWork)

	userGroup := engine.Group("/user")
	userGroup.POST("/douyin",
		timeout.New(
			timeout.WithTimeout(120*time.Second),
			timeout.WithHandler(service.AddDouyinUser),
		))
	userGroup.POST("/douyin/:id/refresh", service.RefreshDouyinUser)
	userGroup.POST("/douyin/:id/manage", service.ManageDouyinUser)
	userGroup.DELETE("/douyin/:id", service.DeleteDouyinUser)
	userGroup.GET("/douyin", service.GetAllDouyinUser)
	userGroup.PUT("/douyin/:id", service.UpdateDouyinUser)
	engine.GET("/label", service.GetAllLabel)
	return engine
}
