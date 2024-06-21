package browser

import (
	"context"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

func Browse(c context.Context, fn func(ctx context.Context, cancel context.CancelFunc) error, options ...chromedp.ExecAllocatorOption) error {
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
