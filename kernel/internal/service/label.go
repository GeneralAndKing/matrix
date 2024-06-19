package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"net/http"
)

func GetAllLabel(c *gin.Context) {
	var (
		labels []model.Label
	)

	err := database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Find(&labels); tx.Error != nil {
			return fmt.Errorf("failed to found labels: %w", tx.Error)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, labels)
}
