package screenshot

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"snapup/internal/models"
)

// Processor 图片处理器接口
type Processor interface {
	Process(data []byte, style models.MockupStyle, background string) ([]byte, error)
}

// ImageProcessor 图片处理器
type ImageProcessor struct{}

// NewImageProcessor 创建图片处理器
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// Process 处理图片，应用样式
func (p *ImageProcessor) Process(data []byte, style models.MockupStyle, bgColor string) ([]byte, error) {
	// 解码原始截图
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	var resultImg image.Image

	switch style {
	case models.StyleNone:
		resultImg = img
	case models.StyleGlass:
		resultImg = p.applyGlassStyle(img, bgColor)
	case models.StyleDevice:
		resultImg = p.applyDeviceStyle(img, bgColor)
	case models.StyleFloating:
		resultImg = p.applyFloatingStyle(img, bgColor)
	default:
		resultImg = img
	}

	// 编码为 PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, resultImg); err != nil {
		return nil, fmt.Errorf("编码图片失败: %w", err)
	}

	return buf.Bytes(), nil
}

// applyGlassStyle 应用玻璃模糊边缘样式
func (p *ImageProcessor) applyGlassStyle(img image.Image, bgColor string) image.Image {
	bounds := img.Bounds()
	padding := 40
	borderRadius := 20

	newWidth := bounds.Dx() + padding*2
	newHeight := bounds.Dy() + padding*2

	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 填充背景
	bg := p.parseColor(bgColor, color.RGBA{R: 240, G: 242, B: 245, A: 255})
	draw.Draw(result, result.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// 绘制半透明边框（模拟玻璃效果）
	borderColor := color.RGBA{R: 255, G: 255, B: 255, A: 100}
	p.drawRoundedRect(result, padding-5, padding-5, bounds.Dx()+10, bounds.Dy()+10, borderRadius+5, borderColor)

	// 绘制原图
	draw.Draw(result, image.Rect(padding, padding, padding+bounds.Dx(), padding+bounds.Dy()), img, bounds.Min, draw.Over)

	return result
}

// applyDeviceStyle 应用设备包裹样式
func (p *ImageProcessor) applyDeviceStyle(img image.Image, bgColor string) image.Image {
	bounds := img.Bounds()
	padding := 60
	deviceBorder := 8

	newWidth := bounds.Dx() + padding*2
	newHeight := bounds.Dy() + padding*2

	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 填充背景
	bg := p.parseColor(bgColor, color.RGBA{R: 240, G: 242, B: 245, A: 255})
	draw.Draw(result, result.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// 绘制设备边框（模拟设备外壳）
	deviceColor := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	p.drawRoundedRect(result, padding-deviceBorder, padding-deviceBorder,
		bounds.Dx()+deviceBorder*2, bounds.Dy()+deviceBorder*2, 15, deviceColor)

	// 绘制屏幕
	draw.Draw(result, image.Rect(padding, padding, padding+bounds.Dx(), padding+bounds.Dy()), img, bounds.Min, draw.Over)

	return result
}

// applyFloatingStyle 应用浮动阴影样式
func (p *ImageProcessor) applyFloatingStyle(img image.Image, bgColor string) image.Image {
	bounds := img.Bounds()
	padding := 50
	shadowOffset := 20

	newWidth := bounds.Dx() + padding*2
	newHeight := bounds.Dy() + padding*2

	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 填充背景
	bg := p.parseColor(bgColor, color.RGBA{R: 240, G: 242, B: 245, A: 255})
	draw.Draw(result, result.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// 绘制阴影
	shadowColor := color.RGBA{R: 0, G: 0, B: 0, A: 30}
	for i := 0; i < shadowOffset; i++ {
		alpha := uint8(30 * (1 - float64(i)/float64(shadowOffset)))
		shadowColor.A = alpha
		p.drawRoundedRect(result, padding+i, padding+i+shadowOffset/2, bounds.Dx(), bounds.Dy(), 10, shadowColor)
	}

	// 绘制原图
	draw.Draw(result, image.Rect(padding, padding, padding+bounds.Dx(), padding+bounds.Dy()), img, bounds.Min, draw.Over)

	return result
}

// drawRoundedRect 绘制圆角矩形
func (p *ImageProcessor) drawRoundedRect(img *image.RGBA, x, y, w, h, r int, c color.Color) {
	// 简化版圆角矩形，使用直角代替（完整实现需要更复杂的算法）
	rect := image.Rect(x, y, x+w, y+h)
	draw.Draw(img, rect, &image.Uniform{c}, image.Point{}, draw.Over)
}

// parseColor 解析颜色字符串
func (p *ImageProcessor) parseColor(colorStr string, defaultColor color.RGBA) color.RGBA {
	if colorStr == "" {
		return defaultColor
	}

	// 简单的十六进制颜色解析
	if len(colorStr) == 7 && colorStr[0] == '#' {
		var r, g, b uint8
		fmt.Sscanf(colorStr, "#%02x%02x%02x", &r, &g, &b)
		return color.RGBA{R: r, G: g, B: b, A: 255}
	}

	// 预定义颜色
	colors := map[string]color.RGBA{
		"white":     {255, 255, 255, 255},
		"black":     {0, 0, 0, 255},
		"gray":      {128, 128, 128, 255},
		"lightgray": {211, 211, 211, 255},
		"blue":      {66, 135, 245, 255},
		"green":     {76, 175, 80, 255},
		"red":       {244, 67, 54, 255},
	}

	if c, ok := colors[colorStr]; ok {
		return c
	}

	return defaultColor
}

// Helper function to calculate distance
func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}
