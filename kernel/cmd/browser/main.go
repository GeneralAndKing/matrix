package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

// Cookie is used to set chromedp_ext cookies.
type Cookie struct {
	Name     string  `json:"name" yaml:"name"`
	Value    string  `json:"value" yaml:"value"`
	Domain   string  `json:"domain" yaml:"domain"`
	Path     string  `json:"path" yaml:"path"`
	Expires  float64 `json:"expires" yaml:"expires"`
	HTTPOnly bool    `json:"httpOnly" yaml:"httpOnly"`
	Secure   bool    `json:"secure" yaml:"secure"`
}

func main() {
	var (
		qrXpath      = `//div[starts-with(@class,"qrcode-image")]/img[1]`
		refreshXpath = `//div[starts-with(@class,"qrcode-image")]/div`
		qrData       string
	)

	// Define options to pass into initial context.
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	// Create a context with options.
	initialCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// Create new context off the initial context.
	ctx, cancel := chromedp.NewContext(initialCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://creator.douyin.com"),
		chromedp.WaitVisible(qrXpath),
		chromedp.AttributeValue(qrXpath, "src", &qrData, nil),
	)
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	chromedp.Run(ctx,
		chromedp.WaitVisible(refreshXpath),
		chromedp.Click(refreshXpath),
	)
	//二次验证
	chromedp.Run(ctx,
		chromedp.WaitVisible(`//div[@class="second-verify-panel"]`),
		chromedp.Click(`//p[text()="接收短信验证"]/..`),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`//p[text()="获取验证码"]/..`),

		chromedp.SendKeys(`//*[@id="uc-second-verify"]//input[@placeholder="请输入验证码"]`, "xxx"),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`//*[@id="uc-second-verify"]//div[text()="验证"]`),
	)
	//保存cookies
	chromedp.Run(ctx, chromedp.WaitVisible(`//*[@id="douyin-creator-master-side-upload"]`))
	if err != nil {
		log.Fatal(err)
	}
}
