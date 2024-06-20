package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/model/dto"
	"kernel/pkg/chromedp_ext"
	"kernel/pkg/message"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func browse(c context.Context, fn func(ctx context.Context, cancel context.CancelFunc) error, options ...chromedp.ExecAllocatorOption) error {
	var (
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.ExecPath("/Users/klein/Projects/matrix/app/dist/electron/Packaged/mac/Quasar App.app/Contents/MacOS/Quasar App"),
			chromedp.Flag("headless", true),
		)
	)
	opts = append(opts, options...)

	// Create a context with options.
	initialCtx, cancel := chromedp.NewExecAllocator(c, opts...)
	defer cancel()
	// Create new context off the initial context.
	chromedpCtx, chromedpCancel := chromedp.NewContext(initialCtx, chromedp.WithLogf(zap.S().Infof))
	return fn(chromedpCtx, chromedpCancel)
}

func refreshDouyinUser(c context.Context, user model.DouyinUser) (name, douyinId, description, avatar string, cookies []chromedp_ext.Cookie, err error) {
	err = browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		return chromedp.Run(ctx,
			chromedp_ext.LoadCookies(user.Cookies),
			chromedp.Navigate("https://creator.douyin.com/creator-micro/home"),
			//等待10秒 如果超时则需要重新登陆
			chromedp_ext.WithTimeOut(10*time.Second,
				chromedp.Tasks{chromedp.WaitVisible(`//*[@id="douyin-creator-master-side-upload"]`)}),
			chromedp_ext.SaveCookies(&cookies),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]`, &name),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[4]`, &description),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[3]`, &douyinId),
			chromedp.AttributeValue(`//*[@id="sub-app"]/div/div[2]/div[1]/div[1]/img`, "src", &avatar, nil),
		)
	})
	if err == nil {
		douyinId = douyinId[strings.Index(douyinId, "：")+3:]
	}
	return

}

func addDouyinUser(c context.Context) (name, douyinId, description, avatar string, cookies []chromedp_ext.Cookie, err error) {
	err = browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		return chromedp.Run(ctx,
			chromedp.Navigate("https://creator.douyin.com"),
			//等待抖音登陆成功
			chromedp.WaitVisible(`//*[@id="douyin-creator-master-side-upload"]`),
			chromedp_ext.SaveCookies(&cookies),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]`, &name),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[4]`, &description),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[3]`, &douyinId),
			chromedp.AttributeValue(`//*[@id="sub-app"]/div/div[2]/div[1]/div[1]/img`, "src", &avatar, nil),
		)
	}, chromedp.Flag("headless", false))
	if err == nil {
		douyinId = douyinId[strings.Index(douyinId, "：")+3:]
	}
	return
}

func manageDouyinUser(c context.Context, user model.DouyinUser) error {
	return browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		_ = message.Fetch(message.WS).Publish(message.Message{
			Type:    0,
			Content: fmt.Sprintf("正在管理 %s 抖音号", user.DouyinId),
		})
		err := chromedp.Run(ctx,
			chromedp_ext.LoadCookies(user.Cookies),
			chromedp.Navigate("https://creator.douyin.com/creator-micro/home"),
			chromedp.Sleep(24*time.Hour),
		)
		if err != nil && errors.Is(context.Canceled, err) {
			return nil
		}
		return err
	}, chromedp.Flag("headless", false))
}

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
	name, douyinId, description, avatar, cookies, err := addDouyinUser(c)
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
		name, douyinId, description, avatar, cookies, refreshErr = refreshDouyinUser(c, user)
		if refreshErr != nil {
			zap.L().Warn("failed to refresh douyin user, re login", zap.Error(refreshErr), zap.Uint("id", user.ID))
			name, douyinId, description, avatar, cookies, loginErr = addDouyinUser(c)
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
		outputs []dto.DouyinUserOutput
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
		return manageDouyinUser(c, user)
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusNoContent)
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
	c.Status(http.StatusNoContent)

}
