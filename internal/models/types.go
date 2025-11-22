package models

// DeviceType 表示设备类型
type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceLaptop  DeviceType = "laptop"
	DeviceTablet  DeviceType = "tablet"
	DeviceMobile  DeviceType = "mobile"
)

// MockupStyle 表示 mockup 样式
type MockupStyle string

const (
	StyleNone     MockupStyle = "none"     // 无样式，纯截图
	StyleGlass    MockupStyle = "glass"    // 玻璃模糊边缘包裹
	StyleDevice   MockupStyle = "device"   // 设备包裹
	StyleFloating MockupStyle = "floating" // 浮动阴影
)

// DeviceConfig 设备配置
type DeviceConfig struct {
	Width  int64
	Height int64
	Scale  float64
	Mobile bool
}

// ScreenshotRequest 截图请求
type ScreenshotRequest struct {
	URL        string      `json:"url"`
	Device     DeviceType  `json:"device"`
	Style      MockupStyle `json:"style"`
	Delay      int         `json:"delay"`      // 延迟时间(毫秒)
	FullPage   bool        `json:"full_page"`  // 是否全页截图
	Quality    int         `json:"quality"`    // 图片质量 (1-100)
	Background string      `json:"background"` // 背景颜色
}

// ScreenshotResponse 截图响应
type ScreenshotResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// GetDeviceConfig 获取设备配置
func GetDeviceConfig(deviceType DeviceType) DeviceConfig {
	configs := map[DeviceType]DeviceConfig{
		DeviceDesktop: {
			Width:  1920,
			Height: 1080,
			Scale:  1.0,
			Mobile: false,
		},
		DeviceLaptop: {
			Width:  1440,
			Height: 900,
			Scale:  1.0,
			Mobile: false,
		},
		DeviceTablet: {
			Width:  768,
			Height: 1024,
			Scale:  2.0,
			Mobile: true,
		},
		DeviceMobile: {
			Width:  375,
			Height: 812,
			Scale:  2.0,
			Mobile: true,
		},
	}

	config, exists := configs[deviceType]
	if !exists {
		return configs[DeviceDesktop]
	}
	return config
}
