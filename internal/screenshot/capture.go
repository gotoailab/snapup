package screenshot

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gotoailab/snapup/internal/models"

	"github.com/chromedp/chromedp"
)

// Capturer 截图捕获器接口
type Capturer interface {
	Capture(ctx context.Context, req models.ScreenshotRequest) ([]byte, error)
}

// ChromeCapture Chrome 截图捕获器
type ChromeCapture struct {
	chromeWSURL string
}

// NewChromeCapture 创建 Chrome 截图捕获器
func NewChromeCapture() *ChromeCapture {
	return &ChromeCapture{
		chromeWSURL: os.Getenv("CHROME_WS_URL"),
	}
}

// Capture 执行截图
func (c *ChromeCapture) Capture(ctx context.Context, req models.ScreenshotRequest) ([]byte, error) {
	// 获取设备配置
	deviceConfig := models.GetDeviceConfig(req.Device)

	var taskCtx context.Context
	var cancel context.CancelFunc

	// 根据是否配置了远程 Chrome 来创建上下文
	if c.chromeWSURL != "" {
		// 使用远程 Chrome（Docker 容器）
		allocCtx, allocCancel := chromedp.NewRemoteAllocator(ctx, c.chromeWSURL)
		defer allocCancel()

		taskCtx, cancel = chromedp.NewContext(allocCtx)
		defer cancel()
	} else {
		// 使用本地 Chrome
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("disable-setuid-sandbox", true),
			chromedp.WindowSize(int(deviceConfig.Width), int(deviceConfig.Height)),
		)

		allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
		defer allocCancel()

		taskCtx, cancel = chromedp.NewContext(allocCtx)
		defer cancel()
	}

	// 设置超时
	taskCtx, cancel = context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel()

	// 截图缓冲区
	var buf []byte

	// 构建任务列表
	tasks := chromedp.Tasks{}
	
	// 设置视口和设备模拟
	if deviceConfig.Mobile {
		// 移动设备模拟
		tasks = append(tasks,
			chromedp.EmulateViewport(deviceConfig.Width, deviceConfig.Height, 
				chromedp.EmulateScale(deviceConfig.Scale)),
		)
	} else {
		// 桌面设备
		tasks = append(tasks,
			chromedp.EmulateViewport(deviceConfig.Width, deviceConfig.Height),
		)
	}
	
	// 导航到目标 URL
	tasks = append(tasks, chromedp.Navigate(req.URL))

	// 如果需要延迟
	if req.Delay > 0 {
		tasks = append(tasks, chromedp.Sleep(time.Duration(req.Delay)*time.Millisecond))
	} else {
		// 默认等待页面加载完成
		tasks = append(tasks, chromedp.Sleep(1*time.Second))
	}

	// 隐藏滚动条
	hideScrollbarJS := `
		(function() {
			const style = document.createElement('style');
			style.textContent = '* { scrollbar-width: none !important; -ms-overflow-style: none !important; } *::-webkit-scrollbar { display: none !important; width: 0 !important; height: 0 !important; } body { overflow: -moz-scrollbars-none !important; }';
			document.head.appendChild(style);
		})();
	`
	tasks = append(tasks, chromedp.Evaluate(hideScrollbarJS, nil))

	// 执行截图
	if req.FullPage {
		tasks = append(tasks, chromedp.FullScreenshot(&buf, 100))
	} else {
		tasks = append(tasks, chromedp.CaptureScreenshot(&buf))
	}

	// 执行任务
	if err := chromedp.Run(taskCtx, tasks); err != nil {
		return nil, fmt.Errorf("截图失败: %w", err)
	}

	return buf, nil
}
