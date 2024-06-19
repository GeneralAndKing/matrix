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
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func refreshDouyinUser(ctx context.Context, id uint) (model.DouyinUser, error) {
	var (
		user model.DouyinUser
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("/Users/klein/Projects/matrix/app/dist/electron/Packaged/mac/Quasar App.app/Contents/MacOS/Quasar App"),
			chromedp.Flag("start-fullscreen", true),
			chromedp.Headless,
		)
		name        string
		douyinId    string
		description string
		avatar      string
		cookies     []chromedp_ext.Cookie
	)
	return user, database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if tx := db.Preload("Labels").Find(&user, id); tx.Error != nil {
			return fmt.Errorf("failed to find douyin user: %w", tx.Error)
		}
		// Create a context with options.
		initialCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
		// Create new context off the initial context.
		chromedpCtx, _ := chromedp.NewContext(initialCtx, chromedp.WithLogf(zap.S().Infof))

		defer cancel()
		err := chromedp.Run(chromedpCtx,
			chromedp_ext.LoadCookies(user.Cookies),
			chromedp.Navigate("https://creator.douyin.com/creator-micro/home"),
			chromedp_ext.WithTimeOut(10*time.Second,
				chromedp.Tasks{chromedp.WaitVisible(`//*[@id="douyin-creator-master-side-upload"]`)}),
			chromedp_ext.SaveCookies(&cookies),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]`, &name),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[4]`, &description),
			chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[3]`, &douyinId),
			chromedp.AttributeValue(`//*[@id="sub-app"]/div/div[2]/div[1]/div[1]/img`, "src", &avatar, nil),
		)
		if err != nil {
			user.Expired = true
			zap.L().Warn("failed to refresh douyin user, set expired to true", zap.Error(err), zap.Uint("id", user.ID))
		} else {
			user.Expired = false
			user.Name = name
			user.Description = description
			user.DouyinId = douyinId[strings.Index(douyinId, "：")+3:]
			user.Avatar = avatar
			user.Cookies = cookies
		}
		if tx := db.Save(&user); tx.Error != nil {
			return fmt.Errorf("failed to save douyin user: %w", tx.Error)
		}
		return nil
	})
}

func AddDouyinUser(c *gin.Context) {

	var (
		input dto.AddDouyinUserInput
		opts  = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("/Users/klein/Projects/matrix/app/dist/electron/Packaged/mac/Quasar App.app/Contents/MacOS/Quasar App"),
			chromedp.Flag("headless", false),
			chromedp.Flag("start-fullscreen", true),
		)
		cookies         []chromedp_ext.Cookie
		name            string
		douyinId        string
		description     string
		avatar          string
		user            model.DouyinUser
		labels          []model.Label
		needSavedLabels []model.Label
	)
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Create a context with options.
	initialCtx, cancel := chromedp.NewExecAllocator(c, opts...)
	// Create new context off the initial context.
	chromedpCtx, _ := chromedp.NewContext(initialCtx, chromedp.WithLogf(zap.S().Infof))

	defer cancel()
	err := chromedp.Run(chromedpCtx,
		chromedp.Navigate("https://creator.douyin.com"),
		//等待抖音登陆成功
		chromedp.WaitVisible(`//*[@id="douyin-creator-master-side-upload"]`),
		chromedp_ext.SaveCookies(&cookies),
		chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[1]/div[1]`, &name),
		chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[4]`, &description),
		chromedp.Text(`//*[@id="sub-app"]/div/div[2]/div[1]/div[2]/div[1]/div[3]`, &douyinId),
		chromedp.AttributeValue(`//*[@id="sub-app"]/div/div[2]/div[1]/div[1]/img`, "src", &avatar, nil),
	)
	if err != nil {
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
				user.DouyinId = douyinId[strings.Index(douyinId, "：")+3:]
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
	c.JSON(http.StatusOK, user.Output())

}

func RefreshDouyinUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("douyin id should be uint type: %w", err))
		return
	}
	user, err := refreshDouyinUser(c, uint(id))
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
	c.JSON(http.StatusOK, user.Output())
}
