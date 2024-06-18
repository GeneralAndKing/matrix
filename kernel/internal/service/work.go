package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/model/dto"
	"kernel/internal/model/page"
	"math"
	"net/http"
)

func GetWorkPage(ctx *gin.Context) {
	var (
		totalElements int64
		totalPages    int64
		works         []model.Work
		input         dto.WorkPageInput
	)
	err := ctx.ShouldBindQuery(&input)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err = database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		db = db.Model(&model.Work{})
		if input.Types != nil {
			db = db.Where("type in ?", input.Types)
		}
		if result := db.Count(&totalElements); result.Error != nil {
			return fmt.Errorf("count works: %w", result.Error)
		}
		totalPages = int64(math.Ceil(float64(totalElements / int64(input.Size))))
		offset := input.Page * input.Size
		if result := db.Offset(offset).Limit(input.Size).Find(&works); result.Error != nil {
			return fmt.Errorf("get page works: %w", result.Error)
		}
		return nil
	}); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, page.FromSliceWithMap(works, input.Page, input.Size, totalPages, totalElements,
		func(i model.Work) dto.WorkOutput {
			return i.ToOutput()
		}))
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
	database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		db.Save(model.Work{
			Type:        input.Type,
			Title:       input.Title,
			Description: input.Description,
			Paths:       input.Paths,
		})
		return nil
	})
}
