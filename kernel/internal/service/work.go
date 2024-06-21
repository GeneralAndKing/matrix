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

func GetAllCreation(ctx *gin.Context) {
	var (
		creations []model.Creation
		outputs   []dto.CreationOutput
	)

	if err := database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if result := db.Find(&creations); result.Error != nil {
			return fmt.Errorf("failed to found creations: %w", result.Error)
		}
		return nil
	}); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	for _, creation := range creations {
		outputs = append(outputs, creation.Output())
	}
	ctx.JSON(http.StatusOK, outputs)
}

func AddCreation(ctx *gin.Context) {
	var (
		input dto.AddCreationInput
	)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	creation := model.Creation{
		Type:        input.Type,
		Title:       input.Title,
		Description: input.Description,
		Paths:       input.Paths,
	}
	err = database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if result := db.Save(&creation); result.Error != nil {
			return fmt.Errorf("failed to save creation: %w", result.Error)
		}
		return nil
	})
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, creation)
}

func PublishCreation(c *gin.Context) {
	var (
		creation model.Creation
		input    dto.PublishCreationInput
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("creation id should be uint type: %w", err))
		return
	}
	if err = c.ShouldBindJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Find(&creation, id); tx.Error != nil {
			return fmt.Errorf("failed to find creation: %w", tx.Error)
		}
		for _, userInput := range input.DouyinUserInputs {
			var douyinUser model.DouyinUser
			if userTx := db.Find(&douyinUser, userInput.ID); userTx.Error != nil {
				return fmt.Errorf("failed to find user %s: %w", userInput.ID, userTx.Error)
			}
			douyinCreation := model.DouyinCreation{
				DouyinUser:        douyinUser,
				Creation:          creation,
				CreationID:        creation.ID,
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
				Status:            enum.PendingCreationStatus,
			}
			if douyinCreationTx := db.Save(&douyinCreation); douyinCreationTx.Error != nil {
				return fmt.Errorf("failed to save douyin creation: %w", douyinCreationTx.Error)
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
