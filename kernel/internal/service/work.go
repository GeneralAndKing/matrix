package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/model/dto"
	"kernel/internal/model/enum"
	"net/http"
	"strconv"
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

func PublishWork(c *gin.Context) {
	var (
		work  model.Work
		input dto.PublishWorkInput
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("work id should be uint type: %w", err))
		return
	}
	if err = c.ShouldBindJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Find(&work, id); tx.Error != nil {
			return fmt.Errorf("failed to find work: %w", tx.Error)
		}
		for _, userInput := range input.DouyinUserInputs {
			var douyinUser model.DouyinUser
			if userTx := db.Find(&douyinUser, userInput.ID); userTx.Error != nil {
				return fmt.Errorf("failed to find user %s: %w", userInput.ID, userTx.Error)
			}
			douyinWork := model.DouyinWork{
				DouyinUser:        douyinUser,
				Work:              work,
				WorkID:            work.ID,
				Title:             userInput.Title,
				Description:       userInput.Description,
				VideoCoverPath:    userInput.VideoCoverPath,
				Location:          userInput.Location,
				Paster:            userInput.Paster,
				CollectionName:    userInput.CollectionName,
				CollectionNum:     userInput.CollectionNum,
				AssociatedHotspot: userInput.AssociatedHotspot,
				SyncToToutiao:     userInput.SyncToToutiao,
				AllowedToSave:     userInput.AllowedToSave,
				WhoCanWatch:       userInput.WhoCanWatch,
				ReleaseTime:       userInput.ReleaseTime,
				Status:            enum.PendingWorkStatus,
			}
			if douyinWorkTx := db.Save(&douyinWork); douyinWorkTx.Error != nil {
				return fmt.Errorf("failed to save douyin work: %w", douyinWorkTx.Error)
			}
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
