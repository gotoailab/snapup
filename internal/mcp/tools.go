package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotoailab/snapup/internal/models"
	"github.com/gotoailab/snapup/internal/screenshot"
)

// ScreenshotToolHandler 截图工具处理器
type ScreenshotToolHandler struct {
	service   *screenshot.Service
	outputDir string
}

// NewScreenshotToolHandler 创建截图工具处理器
func NewScreenshotToolHandler(service *screenshot.Service, outputDir string) *ScreenshotToolHandler {
	return &ScreenshotToolHandler{
		service:   service,
		outputDir: outputDir,
	}
}

// RegisterScreenshotTools 注册所有截图相关的工具
func (h *ScreenshotToolHandler) RegisterScreenshotTools(server *Server) error {
	// 注册截图工具
	screenshotTool := Tool{
		Name:        "take_screenshot",
		Description: "获取指定网站的屏幕截图。支持不同设备尺寸（桌面、笔记本、平板、手机）和样式（无样式、玻璃风格、设备边框、浮动阴影）。返回 base64 编码的 PNG 图片。",
	}

	// 定义工具输入模式
	inputSchema := ToolInput{
		Type: "object",
		Properties: map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "要截图的网站 URL（必须包含 http:// 或 https://）",
			},
			"device": map[string]interface{}{
				"type":        "string",
				"description": "设备类型",
				"enum":        []string{"desktop", "laptop", "tablet", "mobile"},
				"default":     "desktop",
			},
			"style": map[string]interface{}{
				"type":        "string",
				"description": "截图样式",
				"enum":        []string{"none", "glass", "device", "floating"},
				"default":     "none",
			},
			"full_page": map[string]interface{}{
				"type":        "boolean",
				"description": "是否截取全页（整个页面内容）还是仅可见区域",
				"default":     false,
			},
			"delay": map[string]interface{}{
				"type":        "integer",
				"description": "截图前的延迟时间（毫秒），用于等待页面加载完成",
				"default":     1000,
				"minimum":     0,
				"maximum":     30000,
			},
			"quality": map[string]interface{}{
				"type":        "integer",
				"description": "图片质量（1-100）",
				"default":     90,
				"minimum":     1,
				"maximum":     100,
			},
			"background": map[string]interface{}{
				"type":        "string",
				"description": "背景颜色（十六进制格式，如 #f0f2f5，或预定义颜色名称）",
				"default":     "#f0f2f5",
			},
		},
		Required: []string{"url"},
	}

	schemaBytes, err := json.Marshal(inputSchema)
	if err != nil {
		return fmt.Errorf("序列化输入模式失败: %w", err)
	}

	screenshotTool.InputSchema = schemaBytes

	// 注册工具
	server.RegisterTool(screenshotTool, h.handleTakeScreenshot)

	// 注册设备信息工具
	devicesInfoTool := Tool{
		Name:        "get_devices_info",
		Description: "获取所有支持的设备类型及其屏幕尺寸信息",
	}

	devicesInfoSchema := ToolInput{
		Type:       "object",
		Properties: map[string]interface{}{},
		Required:   []string{},
	}

	devicesInfoSchemaBytes, err := json.Marshal(devicesInfoSchema)
	if err != nil {
		return fmt.Errorf("序列化设备信息输入模式失败: %w", err)
	}

	devicesInfoTool.InputSchema = devicesInfoSchemaBytes
	server.RegisterTool(devicesInfoTool, h.handleGetDevicesInfo)

	// 注册样式信息工具
	stylesInfoTool := Tool{
		Name:        "get_styles_info",
		Description: "获取所有支持的截图样式及其描述",
	}

	stylesInfoSchema := ToolInput{
		Type:       "object",
		Properties: map[string]interface{}{},
		Required:   []string{},
	}

	stylesInfoSchemaBytes, err := json.Marshal(stylesInfoSchema)
	if err != nil {
		return fmt.Errorf("序列化样式信息输入模式失败: %w", err)
	}

	stylesInfoTool.InputSchema = stylesInfoSchemaBytes
	server.RegisterTool(stylesInfoTool, h.handleGetStylesInfo)

	return nil
}

// handleTakeScreenshot 处理截图请求
func (h *ScreenshotToolHandler) handleTakeScreenshot(ctx context.Context, arguments map[string]interface{}) (*CallToolResult, error) {
	// 解析参数
	url, _ := arguments["url"].(string)
	if url == "" {
		return &CallToolResult{
			Content: []Content{{
				Type: "text",
				Text: "错误：URL 参数是必需的",
			}},
			IsError: true,
		}, nil
	}

	device, _ := arguments["device"].(string)
	if device == "" {
		device = "desktop"
	}

	style, _ := arguments["style"].(string)
	if style == "" {
		style = "none"
	}

	fullPage, _ := arguments["full_page"].(bool)

	delay := 1000
	if d, ok := arguments["delay"].(float64); ok {
		delay = int(d)
	}

	quality := 90
	if q, ok := arguments["quality"].(float64); ok {
		quality = int(q)
	}

	background, _ := arguments["background"].(string)
	if background == "" {
		background = "#f0f2f5"
	}

	// 创建截图请求
	req := models.ScreenshotRequest{
		URL:        url,
		Device:     models.DeviceType(device),
		Style:      models.MockupStyle(style),
		Delay:      delay,
		FullPage:   fullPage,
		Quality:    quality,
		Background: background,
	}

	// 执行截图
	resp, err := h.service.TakeScreenshot(ctx, req)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("截图失败: %v", err),
			}},
			IsError: true,
		}, nil
	}

	if !resp.Success {
		return &CallToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("截图失败: %s", resp.Message),
			}},
			IsError: true,
		}, nil
	}

	// 读取截图文件并转换为 base64
	screenshotPath := filepath.Join(h.outputDir, resp.Filename)
	imageData, err := os.ReadFile(screenshotPath)
	if err != nil {
		return &CallToolResult{
			Content: []Content{{
				Type: "text",
				Text: fmt.Sprintf("读取截图文件失败: %v", err),
			}},
			IsError: true,
		}, nil
	}

	// 转换为 base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// 获取设备配置信息
	deviceConfig := models.GetDeviceConfig(req.Device)

	// 返回结果
	resultText := fmt.Sprintf(`截图成功！

URL: %s
设备: %s (%dx%d)
样式: %s
全页截图: %v
延迟: %d 毫秒
质量: %d%%
文件名: %s

图片已生成为 base64 编码的 PNG 格式。`,
		url, device, deviceConfig.Width, deviceConfig.Height,
		style, fullPage, delay, quality, resp.Filename)

	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: resultText,
			},
			{
				Type:     "image",
				Data:     base64Image,
				MimeType: "image/png",
			},
		},
		IsError: false,
	}, nil
}

// handleGetDevicesInfo 处理获取设备信息请求
func (h *ScreenshotToolHandler) handleGetDevicesInfo(ctx context.Context, arguments map[string]interface{}) (*CallToolResult, error) {
	devices := []struct {
		Type   string
		Name   string
		Width  int64
		Height int64
		Mobile bool
	}{
		{
			Type:   "desktop",
			Name:   "桌面",
			Width:  1920,
			Height: 1080,
			Mobile: false,
		},
		{
			Type:   "laptop",
			Name:   "笔记本",
			Width:  1440,
			Height: 900,
			Mobile: false,
		},
		{
			Type:   "tablet",
			Name:   "平板",
			Width:  768,
			Height: 1024,
			Mobile: true,
		},
		{
			Type:   "mobile",
			Name:   "手机",
			Width:  375,
			Height: 812,
			Mobile: true,
		},
	}

	resultText := "支持的设备类型：\n\n"
	for _, device := range devices {
		resultText += fmt.Sprintf("- %s (%s)\n", device.Name, device.Type)
		resultText += fmt.Sprintf("  尺寸: %dx%d\n", device.Width, device.Height)
		resultText += fmt.Sprintf("  移动设备: %v\n\n", device.Mobile)
	}

	return &CallToolResult{
		Content: []Content{{
			Type: "text",
			Text: resultText,
		}},
		IsError: false,
	}, nil
}

// handleGetStylesInfo 处理获取样式信息请求
func (h *ScreenshotToolHandler) handleGetStylesInfo(ctx context.Context, arguments map[string]interface{}) (*CallToolResult, error) {
	styles := []struct {
		Type        string
		Name        string
		Description string
	}{
		{
			Type:        "none",
			Name:        "无样式",
			Description: "纯净截图，无任何装饰效果",
		},
		{
			Type:        "glass",
			Name:        "玻璃风格",
			Description: "半透明玻璃边框效果，现代感十足",
		},
		{
			Type:        "device",
			Name:        "设备边框",
			Description: "模拟设备外壳边框，看起来像真实设备",
		},
		{
			Type:        "floating",
			Name:        "浮动阴影",
			Description: "悬浮效果带阴影，增加立体感",
		},
	}

	resultText := "支持的截图样式：\n\n"
	for _, style := range styles {
		resultText += fmt.Sprintf("- %s (%s)\n", style.Name, style.Type)
		resultText += fmt.Sprintf("  %s\n\n", style.Description)
	}

	return &CallToolResult{
		Content: []Content{{
			Type: "text",
			Text: resultText,
		}},
		IsError: false,
	}, nil
}
