package browser

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"kernel/internal/model"
	"kernel/pkg/chromedp_ext"
	"strings"
	"time"
)

var (
	uploadButtonXpath = `//*[@id="douyin-creator-master-side-upload"]`
)

func RefreshDouyinUser(c context.Context, user model.DouyinUser) (name, douyinId, description, avatar string, cookies []chromedp_ext.Cookie, err error) {
	err = Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		return chromedp.Run(ctx,
			chromedp_ext.LoadCookies(user.Cookies),
			chromedp.Navigate("https://creator.douyin.com/creator-micro/home"),
			//等待10秒 如果超时则需要重新登陆
			chromedp_ext.WithTimeOut(10*time.Second,
				chromedp.Tasks{chromedp.WaitVisible(uploadButtonXpath)}),
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

func AddDouyinUser(c context.Context) (name, douyinId, description, avatar string, cookies []chromedp_ext.Cookie, err error) {
	err = Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		return chromedp.Run(ctx,
			chromedp.Navigate("https://creator.douyin.com"),
			//等待抖音登陆成功
			chromedp.WaitVisible(uploadButtonXpath),
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

func ManageDouyinUser(c context.Context, user model.DouyinUser) error {
	return Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
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

func PublishDouyinCreation(c context.Context, douyinCreation model.DouyinCreation) error {
	return Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		err := chromedp.Run(ctx,
			chromedp_ext.LoadCookies(douyinCreation.DouyinUser.Cookies),
			chromedp.Navigate("https://creator.douyin.com/creator-micro/home"),
			chromedp.WaitVisible(uploadButtonXpath),
			chromedp.Click(uploadButtonXpath),
			chromedp.SetUploadFiles(`//input[@name="upload-btn"]`, douyinCreation.Creation.Paths),
			chromedp_ext.WithTimeOut(15*time.Second, chromedp.Tasks{chromedp.WaitVisible(`//div[text()="发布视频"]`)}),
			chromedp_ext.SendKeys(`//input[@placeholder="好的作品标题可获得更多浏览"]`, douyinCreation.Title),
			chromedp_ext.SendKeys(`//div[@data-placeholder="添加作品简介"]`, douyinCreation.Description),
			chromedp.Click(`//div[text()="选择封面"]`),
			chromedp.Click(`//div[text()="上传封面"]`),
			chromedp.SetUploadFiles(`//input[@class="semi-upload-hidden-input"]`, []string{douyinCreation.VideoCoverPath}),
		)
		if err != nil && errors.Is(context.Canceled, err) {
			return nil
		}
		return err
	})
}
