package screenshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotoailab/snapup/internal/models"

	"github.com/google/uuid"
)

// Service 截图服务
type Service struct {
	capturer  Capturer
	outputDir string
}

// NewService 创建截图服务
func NewService(outputDir string) *Service {
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("创建输出目录失败: %v", err))
	}

	return &Service{
		capturer:  NewChromeCapture(),
		outputDir: outputDir,
	}
}

// TakeScreenshot 执行截图
func (s *Service) TakeScreenshot(ctx context.Context, req models.ScreenshotRequest) (*models.ScreenshotResponse, error) {
	// 验证请求
	if err := s.validateRequest(&req); err != nil {
		return &models.ScreenshotResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 执行截图
	data, err := s.capturer.Capture(ctx, req)
	if err != nil {
		return &models.ScreenshotResponse{
			Success: false,
			Message: fmt.Sprintf("截图失败: %v", err),
		}, nil
	}

	// 保存文件（不再处理样式，直接保存原始截图）
	filename := s.generateFilename(req)
	filepath := filepath.Join(s.outputDir, filename)

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return &models.ScreenshotResponse{
			Success: false,
			Message: fmt.Sprintf("保存文件失败: %v", err),
		}, nil
	}

	return &models.ScreenshotResponse{
		Success:  true,
		Message:  "截图成功",
		ImageURL: "/screenshots/" + filename,
		Filename: filename,
	}, nil
}

// validateRequest 验证请求
func (s *Service) validateRequest(req *models.ScreenshotRequest) error {
	if req.URL == "" {
		return fmt.Errorf("URL 不能为空")
	}

	// 设置默认值
	if req.Device == "" {
		req.Device = models.DeviceDesktop
	}
	if req.Style == "" {
		req.Style = models.StyleNone
	}
	if req.Quality == 0 {
		req.Quality = 90
	}
	if req.Background == "" {
		req.Background = "#f0f2f5"
	}

	return nil
}

// generateFilename 生成文件名
func (s *Service) generateFilename(req models.ScreenshotRequest) string {
	id := uuid.New().String()
	return fmt.Sprintf("screenshot_%s_%s.png", req.Device, id)
}

// CleanupOldScreenshots 清理旧截图（可选功能）
func (s *Service) CleanupOldScreenshots(maxFiles int) error {
	files, err := os.ReadDir(s.outputDir)
	if err != nil {
		return err
	}

	if len(files) <= maxFiles {
		return nil
	}

	// 删除最旧的文件
	for i := 0; i < len(files)-maxFiles; i++ {
		filepath := filepath.Join(s.outputDir, files[i].Name())
		os.Remove(filepath)
	}

	return nil
}
