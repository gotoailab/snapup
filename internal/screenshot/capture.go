package screenshot

import (
	"context"
	"fmt"
	"snapup/internal/models"
	"time"

	"github.com/chromedp/chromedp"
)

// Capturer 截图捕获器接口
type Capturer interface {
	Capture(ctx context.Context, req models.ScreenshotRequest) ([]byte, error)
}

// ChromeCapture Chrome 截图捕获器
type ChromeCapture struct{}

// NewChromeCapture 创建 Chrome 截图捕获器
func NewChromeCapture() *ChromeCapture {
	return &ChromeCapture{}
}

// Capture 执行截图
func (c *ChromeCapture) Capture(ctx context.Context, req models.ScreenshotRequest) ([]byte, error) {
	// 获取设备配置
	deviceConfig := models.GetDeviceConfig(req.Device)

	// 创建 Chrome 上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.WindowSize(int(deviceConfig.Width), int(deviceConfig.Height)),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 设置超时
	taskCtx, cancel = context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel()

	// 截图缓冲区
	var buf []byte

	// 构建任务列表
	tasks := chromedp.Tasks{
		chromedp.EmulateViewport(deviceConfig.Width, deviceConfig.Height),
		chromedp.Navigate(req.URL),
	}

	// 如果需要延迟
	if req.Delay > 0 {
		tasks = append(tasks, chromedp.Sleep(time.Duration(req.Delay)*time.Millisecond))
	} else {
		// 默认等待页面加载完成
		tasks = append(tasks, chromedp.Sleep(1*time.Second))
	}

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
