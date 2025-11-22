package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"snapup/internal/models"
	"snapup/internal/screenshot"
)

// Handler HTTP 处理器
type Handler struct {
	screenshotService *screenshot.Service
}

// NewHandler 创建处理器
func NewHandler(screenshotService *screenshot.Service) *Handler {
	return &Handler{
		screenshotService: screenshotService,
	}
}

// HandleScreenshot 处理截图请求
func (h *Handler) HandleScreenshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求
	var req models.ScreenshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	log.Printf("收到截图请求: URL=%s, Device=%s, Style=%s", req.URL, req.Device, req.Style)

	// 执行截图
	resp, err := h.screenshotService.TakeScreenshot(r.Context(), req)
	if err != nil {
		log.Printf("截图失败: %v", err)
		h.sendJSONError(w, fmt.Sprintf("截图失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回响应
	h.sendJSON(w, resp, http.StatusOK)
}

// HandleDevices 获取支持的设备列表
func (h *Handler) HandleDevices(w http.ResponseWriter, r *http.Request) {
	devices := []map[string]interface{}{
		{
			"type":   string(models.DeviceDesktop),
			"name":   "桌面",
			"width":  1920,
			"height": 1080,
		},
		{
			"type":   string(models.DeviceLaptop),
			"name":   "笔记本",
			"width":  1440,
			"height": 900,
		},
		{
			"type":   string(models.DeviceTablet),
			"name":   "平板",
			"width":  768,
			"height": 1024,
		},
		{
			"type":   string(models.DeviceMobile),
			"name":   "手机",
			"width":  375,
			"height": 812,
		},
	}

	h.sendJSON(w, map[string]interface{}{
		"devices": devices,
	}, http.StatusOK)
}

// HandleStyles 获取支持的样式列表
func (h *Handler) HandleStyles(w http.ResponseWriter, r *http.Request) {
	styles := []map[string]interface{}{
		{
			"type":        string(models.StyleNone),
			"name":        "无样式",
			"description": "纯净截图，无任何装饰",
		},
		{
			"type":        string(models.StyleGlass),
			"name":        "玻璃风格",
			"description": "半透明玻璃边框效果",
		},
		{
			"type":        string(models.StyleDevice),
			"name":        "设备边框",
			"description": "模拟设备外壳边框",
		},
		{
			"type":        string(models.StyleFloating),
			"name":        "浮动阴影",
			"description": "悬浮效果带阴影",
		},
	}

	h.sendJSON(w, map[string]interface{}{
		"styles": styles,
	}, http.StatusOK)
}

// HandleHealth 健康检查
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	h.sendJSON(w, map[string]interface{}{
		"status":  "ok",
		"service": "snapup",
	}, http.StatusOK)
}

// sendJSON 发送 JSON 响应
func (h *Handler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendJSONError 发送 JSON 错误响应
func (h *Handler) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	h.sendJSON(w, map[string]interface{}{
		"success": false,
		"message": message,
	}, statusCode)
}
