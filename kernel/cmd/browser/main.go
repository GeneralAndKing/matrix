package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"kernel/pkg/chromedp_ext"
	"kernel/pkg/external_api/douyin"
	"log"
	"os"
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
	douyin.FetchChallengeSug("做法")
	var (
	//qrXpath = `//div[starts-with(@class,"qrcode-image")]/img[1]`
	//refreshXpath = `//div[starts-with(@class,"qrcode-image")]/div`
	//qrData string
	)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/Users/klein/Projects/matrix/app/dist/electron/Packaged/mac/Quasar App.app/Contents/MacOS/Quasar App"),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("kernel", true),
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"),
	)
	// Create a context with options.
	initialCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// Create new context off the initial context.
	ctx, cancel := chromedp.NewContext(initialCtx, chromedp.WithLogf(log.Printf))
	var buf []byte
	defer cancel()
	var test string
	err := chromedp.Run(ctx,
		chromedp_ext.StealthBypass(),
		chromedp.Navigate("https://dappradar.com"),
		chromedp.Sleep(5*time.Second),
		//chromedp.WaitVisible(`//*[@id="root"]/div[2]/div/div[2]/div/section/div[1]/div/div[3]/a[1]`),
		//chromedp.Text(`//*[@id="root"]/div[2]/div/div[2]/div/section/div[1]/div/div[3]/a[1]`, &test),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		panic(err)
	}
	println(test)
	// 将截图保存到文件
	err = os.WriteFile("screenshot.png", buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//ctx, cancel = context.WithCancel(ctx)
	//defer cancel()
	//chromedp.Run(ctx,
	//	chromedp.WaitVisible(refreshXpath),
	//	chromedp.Click(refreshXpath),
	//)
	////二次验证
	//chromedp.Run(ctx,
	//	chromedp.WaitVisible(`//div[@class="second-verify-panel"]`),
	//	chromedp.Click(`//p[text()="接收短信验证"]/..`),
	//	chromedp.Sleep(1*time.Second),
	//	chromedp.Click(`//p[text()="获取验证码"]/..`),
	//
	//	chromedp.SendKeys(`//*[@id="uc-second-verify"]//input[@placeholder="请输入验证码"]`, "xxx"),
	//	chromedp.Sleep(1*time.Second),
	//	chromedp.Click(`//*[@id="uc-second-verify"]//div[text()="验证"]`),
	//)
	////保存cookies
	//chromedp.Run(ctx, chromedp.WaitVisible(`//*[@id="douyin-creator-master-side-upload"]`))
	//if err != nil {
	//	log.Fatal(err)
	//}
}
