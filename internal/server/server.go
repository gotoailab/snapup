package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"snapup/internal/screenshot"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

// Server Web 服务器
type Server struct {
	handler *Handler
	port    int
}

// NewServer 创建服务器
func NewServer(screenshotService *screenshot.Service, port int) *Server {
	return &Server{
		handler: NewHandler(screenshotService),
		port:    port,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 静态文件服务
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("获取静态文件失败: %w", err)
	}

	// 主页
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})

	// 静态资源
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API 路由
	mux.HandleFunc("/api/screenshot", s.handler.HandleScreenshot)
	mux.HandleFunc("/api/devices", s.handler.HandleDevices)
	mux.HandleFunc("/api/styles", s.handler.HandleStyles)
	mux.HandleFunc("/api/health", s.handler.HandleHealth)

	// 截图文件服务
	mux.Handle("/screenshots/", http.StripPrefix("/screenshots/", http.FileServer(http.Dir("./screenshots"))))

	// 应用中间件
	handler := RecoveryMiddleware(LoggingMiddleware(CORSMiddleware(mux)))

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("服务器启动在端口 %d", s.port)
	log.Printf("访问 http://localhost:%d 使用截图服务", s.port)

	return server.ListenAndServe()
}
