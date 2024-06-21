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

func RefreshDouyinUser(c context.Context, user model.DouyinUser) (name, douyinId, description, avatar string, cookies []chromedp_ext.Cookie, err error) {
	err = Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
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

func AddDouyinUser(c context.Context) (name, douyinId, description, avatar string, cookies []chromedp_ext.Cookie, err error) {
	err = Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
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

func PublishDouyinWork(c context.Context, work model.DouyinWork) error {
	return Browse(c, func(ctx context.Context, cancel context.CancelFunc) error {
		err := chromedp.Run(ctx,
			chromedp_ext.LoadCookies(work.DouyinUser.Cookies),
			chromedp.Navigate("https://creator.douyin.com/creator-micro/home"),
		)
		if err != nil && errors.Is(context.Canceled, err) {
			return nil
		}
		return err
	})
}
