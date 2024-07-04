package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
)

func AddDouyinUser(c *gin.Context) {

	var (
		input           dto.AddDouyinUserInput
		user            model.DouyinUser
		labels          []model.Label
		needSavedLabels []model.Label
	)
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	name, douyinId, description, avatar, cookies, err := browser.AddDouyinUser(c)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to add douyin user: %w", err))
		return
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if result := db.Where("douyin_id = ?", douyinId).Preload("Labels").First(&user); result.Error != nil {
			if errors.Is(gorm.ErrRecordNotFound, result.Error) {
				if tx := db.Where("name in ?", input).Find(&labels); tx.Error != nil {
					return fmt.Errorf("failed to found labels: %w", tx.Error)
				}
				for _, labelName := range input {
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
				user.Labels = slices.Concat(labels, needSavedLabels)
				user.Name = name
				user.Description = description
				user.DouyinId = douyinId
				user.Cookies = cookies
				user.Avatar = avatar
				user.Expired = false
				if saveResult := db.Save(&user); saveResult.Error != nil {
					return fmt.Errorf("failed to save douyin user %s: %w", douyinId, saveResult.Error)
				}
			} else {
				return fmt.Errorf("failed to found douyin user %s: %w", douyinId, result.Error)
			}
		} else {
			return fmt.Errorf("douyin id %s already exist", douyinId)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, user.Output())

}

func RefreshDouyinUser(c *gin.Context) {
	var (
		user                                model.DouyinUser
		name, douyinId, description, avatar string
		cookies                             []chromedp_ext.Cookie
		refreshErr, loginErr                error
	)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin id should be uint type: %w", err))
		return
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Preload("Labels").Find(&user, id); tx.Error != nil {
			return fmt.Errorf("failed to find douyin user: %w", tx.Error)
		}
		name, douyinId, description, avatar, cookies, refreshErr = browser.RefreshDouyinUser(c, user)
		if refreshErr != nil {
			zap.L().Warn("failed to refresh douyin user, re login", zap.Error(refreshErr), zap.Uint("id", user.ID))
			name, douyinId, description, avatar, cookies, loginErr = browser.AddDouyinUser(c)
			if loginErr != nil {
				finalErr := errors.Join(refreshErr, loginErr)
				zap.L().Warn("failed to login douyin user", zap.Error(finalErr), zap.Uint("id", user.ID))
				return fmt.Errorf("failed to refresh and relogin douyin user: %w", finalErr)
			}
		}
		user.Name = name
		user.Description = description
		user.DouyinId = douyinId[strings.Index(douyinId, "：")+3:]
		user.Avatar = avatar
		user.Cookies = cookies
		if tx := db.Save(&user); tx.Error != nil {
			return fmt.Errorf("failed to save douyin user: %w", tx.Error)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user.Output())

}

func GetAllDouyinUser(c *gin.Context) {
	var (
		users   []model.DouyinUser
		outputs = make([]dto.DouyinUserOutput, 0)
	)

	err := database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Preload("Labels").Find(&users); tx.Error != nil {
			return fmt.Errorf("failed to found users: %w", tx.Error)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, user := range users {
		outputs = append(outputs, user.Output())
	}
	c.JSON(http.StatusOK, outputs)
}

func UpdateDouyinUser(c *gin.Context) {
	var (
		input           dto.AddDouyinUserInput
		labels          []model.Label
		needSavedLabels []model.Label
		user            model.DouyinUser
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin id should be uint type: %w", err))
		return
	}
	if err = c.ShouldBindBodyWithJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Where("name in ?", input).Find(&labels); tx.Error != nil {
			return fmt.Errorf("failed to found labels: %w", tx.Error)
		}
		if tx := db.Preload("Labels").Find(&user, id); tx.Error != nil {
			return fmt.Errorf("failed to find douyin user: %w", tx.Error)
		}
		for _, labelName := range input {
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
		if err = db.Model(&user).Association("Labels").Replace(slices.Concat(labels, needSavedLabels)); err != nil {
			return fmt.Errorf("failed to Replace douyin user labes: %w", err)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, user.Output())
}

func ManageDouyinUser(c *gin.Context) {
	var (
		user model.DouyinUser
	)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin id should be uint type: %w", err))
		return
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Find(&user, id); tx.Error != nil {
			return fmt.Errorf("failed to found douyin user: %w", tx.Error)
		}
		return browser.ManageDouyinUser(c, user)
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

func DeleteDouyinUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin id should be uint type: %w", err))
		return
	}
	err = database.Sqlite3Transaction(c, func(db *gorm.DB) error {
		if tx := db.Delete(&model.DouyinUser{}, id); tx.Error != nil {
			return fmt.Errorf("failed to delete douyin user: %w", err)
		}
		return nil
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)

}
