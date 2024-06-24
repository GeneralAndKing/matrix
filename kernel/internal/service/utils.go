package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/pkg/external_api/douyin"
	"net/http"
)

func DouyinActivity(c *gin.Context) {
	var douyinUser model.DouyinUser
	err := database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Where("expired = 0").Order("RANDOM()").First(&douyinUser); tx.Error != nil {
			return fmt.Errorf("failed to get douyinUser: %w", tx.Error)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("no logined douyin accounts: %w", err))
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	responses, err := douyin.FetchActivity(douyinUser.Cookies.HttpCookies())
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, responses)
}

func DouyinHotspot(c *gin.Context) {
	var (
		responses []douyin.HotspotResponse
		err       error
	)
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		responses, err = douyin.FetchRecommendHotspot()
	} else {
		responses, err = douyin.FetchSearchHotspot(keyword)
	}
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, responses)
}

func DouyinChallengeSug(c *gin.Context) {
	var (
		responses []douyin.ChallengeSugResponse
		err       error
	)
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		keyword = "热点"
	}
	responses, err = douyin.FetchChallengeSug(keyword)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, responses)
}
