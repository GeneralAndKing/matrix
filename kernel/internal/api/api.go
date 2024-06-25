package api

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kernel/internal/service"
	"net/http"
	"time"
)

func API(ctx context.Context, debug bool) http.Handler {
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
	engine.GET("/message", service.Message)
	engine.GET("/business", service.Business)
	creationGroup := engine.Group("/creation")
	creationGroup.GET("", service.GetAllCreation)
	creationGroup.POST("", service.AddCreation)
	creationGroup.PUT("/:id", service.UpdateCreation)
	creationGroup.GET("/:id", service.GetCreation)
	creationGroup.GET("/douyin/:id", service.GetDouyinCreation)
	creationGroup.POST("/publish", service.PublishCreation)

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

	utilsGroup := engine.Group("/utils")
	utilsGroup.GET("/douyin/activity", service.DouyinActivity)
	utilsGroup.GET("/douyin/hotspot", service.DouyinHotspot)
	utilsGroup.GET("/douyin/challenge", service.DouyinChallengeSug)
	utilsGroup.GET("/douyin/flashmob", service.DouyinFlashmob)

	return engine
}
