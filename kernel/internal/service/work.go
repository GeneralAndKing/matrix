package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/model/dto"
	"net/http"
)

func GetAllWork(ctx *gin.Context) {
	var (
		works   []model.Work
		outputs []dto.WorkOutput
	)

	if err := database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if result := db.Find(&works); result.Error != nil {
			return fmt.Errorf("failed to found works: %w", result.Error)
		}
		return nil
	}); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	for _, work := range works {
		outputs = append(outputs, work.Output())
	}
	ctx.JSON(http.StatusOK, outputs)
}

func AddWork(ctx *gin.Context) {
	var (
		input dto.AddWorkInput
	)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	work := model.Work{
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		Paths:       input.Paths,
	}
	err = database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if result := db.Save(&work); result.Error != nil {
			return fmt.Errorf("failed to save work: %w", result.Error)
		}
		return nil
	})
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, work)
}
