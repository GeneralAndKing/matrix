package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kernel/internal/browser"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/model/dto"
	"kernel/pkg/chromedp_ext"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func GetAllDouyinLive(c *gin.Context) {
	var (
		lives   []model.DouyinLive
		outputs = make([]dto.DouyinLiveOutput, 0)
	)
	err := database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Find(&lives); tx.Error != nil {
			return fmt.Errorf("failed to found douyin live: %w", tx.Error)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, live := range lives {
		outputs = append(outputs, live.Output())
	}
	c.JSON(http.StatusOK, outputs)
}

func DeleteDouyinLive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin live id should be uint type: %w", err))
		return
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Delete(&model.DouyinLive{}, id); tx.Error != nil {
			return fmt.Errorf("failed to delete douyin live: %w", err)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func AddDouyinLive(c *gin.Context) {
	var (
		input           dto.AddDouyinLiveInput
		douyinLive      model.DouyinLive
		labels          []model.Label
		needSavedLabels []model.Label
		selfUrl         string
	)
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := browser.Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		err := chromedp.Run(ctx,

			chromedp_ext.WithTimeOut(10*time.Second, chromedp.Tasks{
				chromedp.Navigate(fmt.Sprintf("https://live.douyin.com/%s", input.LiveId)),
				chromedp.AttributeValue(`//a[@data-e2e="live-room-nickname"]`, "href", &selfUrl, nil),
			}),
		)
		if err != nil {
			return err
		}
		return chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https:%s", selfUrl)),
			chromedp_ext.WithTimeOut(10*time.Second, chromedp.Tasks{
				chromedp.Text(`//div[@data-e2e="user-info"]/div[1]/h1/span`, &douyinLive.Name),
				chromedp.Text(`//div[@data-e2e="user-info"]/p/span`, &douyinLive.DouyinId),
				chromedp.AttributeValue(`//span[contains(@class, "semi-avatar")]/img`, "src", &douyinLive.Avatar, nil),
			},
			))
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to access douyin live: %w", err))
	}
	douyinLive.DouyinId = douyinLive.DouyinId[strings.Index(douyinLive.DouyinId, "：")+3:]
	douyinLive.Avatar = "https:" + douyinLive.Avatar
	douyinLive.LiveId = input.LiveId
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if result := db.Where("live_id = ?", douyinLive.LiveId).Preload("Labels").First(&model.DouyinLive{}); result.Error != nil {
			if errors.Is(gorm.ErrRecordNotFound, result.Error) {
				if tx := db.Where("name in ?", input.Labels).Find(&labels); tx.Error != nil {
					return fmt.Errorf("failed to found labels: %w", tx.Error)
				}
				for _, labelName := range input.Labels {
					if !slices.ContainsFunc(labels, func(label model.Label) bool {
						return label.Name == labelName
					}) {
						needSavedLabels = append(needSavedLabels, model.Label{
							Name: labelName,
						})
					}
				}
				if needSavedLabels != nil {
					if tx := db.Save(&needSavedLabels); tx.Error != nil {
						return fmt.Errorf("failed to save labels: %w", tx.Error)
					}
				}
				douyinLive.Labels = slices.Concat(labels, needSavedLabels)
				if saveResult := db.Save(&douyinLive); saveResult.Error != nil {
					return fmt.Errorf("failed to save douyin live %s: %w", douyinLive.LiveId, saveResult.Error)
				}
			} else {
				return fmt.Errorf("failed to found douyin user %s: %w", douyinLive.LiveId, result.Error)
			}
		} else {
			return fmt.Errorf("live id %s already exist", douyinLive.LiveId)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save douyin live: %w", err))
		return
	}
	c.JSON(http.StatusCreated, douyinLive.Output())
}

func UpdateDouyinLive(c *gin.Context) {
	var (
		input           dto.AddDouyinLiveInput
		labels          []model.Label
		needSavedLabels []model.Label
		douyinLive      model.DouyinLive
	)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin live id should be uint type: %w", err))
		return
	}

	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Where("name in ?", input).Find(&labels); tx.Error != nil {
			return fmt.Errorf("failed to found labels: %w", tx.Error)
		}
		if tx := db.Preload("Labels").Find(&douyinLive, id); tx.Error != nil {
			return fmt.Errorf("failed to find douyin user: %w", tx.Error)
		}
		for _, labelName := range input.Labels {
			if !slices.ContainsFunc(labels, func(label model.Label) bool {
				return label.Name == labelName
			}) {
				needSavedLabels = append(needSavedLabels, model.Label{
					Name: labelName,
				})
			}
		}
		if needSavedLabels != nil {
			if tx := db.Save(&needSavedLabels); tx.Error != nil {
				return fmt.Errorf("failed to save labels: %w", tx.Error)
			}
		}
		if err = db.Model(&douyinLive).Association("Labels").Replace(slices.Concat(labels, needSavedLabels)); err != nil {
			return fmt.Errorf("failed to Replace douyin live labes: %w", err)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save douyin live: %w", err))
		return
	}
	c.JSON(http.StatusCreated, douyinLive.Output())
}
