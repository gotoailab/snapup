package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Server MCP 服务器
type Server struct {
	info         ServerInfo
	capabilities ServerCapabilities
	tools        []Tool
	resources    []Resource
	prompts      []Prompt

	// 工具处理器
	toolHandlers map[string]ToolHandler

	// 资源处理器
	resourceHandlers map[string]ResourceHandler

	// 提示处理器
	promptHandlers map[string]PromptHandler

	mu sync.RWMutex
}

// ToolHandler 工具处理器
type ToolHandler func(ctx context.Context, arguments map[string]interface{}) (*CallToolResult, error)

// ResourceHandler 资源处理器
type ResourceHandler func(ctx context.Context, uri string) (*ReadResourceResult, error)

// PromptHandler 提示处理器
type PromptHandler func(ctx context.Context, arguments map[string]interface{}) (*GetPromptResult, error)

// NewServer 创建 MCP 服务器
func NewServer(name, version string) *Server {
	return &Server{
		info: ServerInfo{
			Name:    name,
			Version: version,
		},
		capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe: false,
			},
			Prompts: &PromptsCapability{
				ListChanged: false,
			},
		},
		tools:            make([]Tool, 0),
		resources:        make([]Resource, 0),
		prompts:          make([]Prompt, 0),
		toolHandlers:     make(map[string]ToolHandler),
		resourceHandlers: make(map[string]ResourceHandler),
		promptHandlers:   make(map[string]PromptHandler),
	}
}

// RegisterTool 注册工具
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tools = append(s.tools, tool)
	s.toolHandlers[tool.Name] = handler
}

// RegisterResource 注册资源
func (s *Server) RegisterResource(resource Resource, handler ResourceHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resources = append(s.resources, resource)
	s.resourceHandlers[resource.URI] = handler
}

// RegisterPrompt 注册提示
func (s *Server) RegisterPrompt(prompt Prompt, handler PromptHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.prompts = append(s.prompts, prompt)
	s.promptHandlers[prompt.Name] = handler
}

// Run 运行服务器（通过 stdio）
func (s *Server) Run(ctx context.Context) error {
	log.Println("MCP Server 启动，使用 stdio 传输")

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 读取一行 JSON-RPC 请求
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("读取请求失败: %w", err)
			}

			// 处理请求
			response := s.handleRequest(ctx, line)

			// 写入响应
			responseBytes, err := json.Marshal(response)
			if err != nil {
				log.Printf("序列化响应失败: %v", err)
				continue
			}

			// 写入响应和换行符
			if _, err := writer.Write(responseBytes); err != nil {
				return fmt.Errorf("写入响应失败: %w", err)
			}
			if err := writer.WriteByte('\n'); err != nil {
				return fmt.Errorf("写入换行符失败: %w", err)
			}
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("刷新输出失败: %w", err)
			}
		}
	}
}

// handleRequest 处理请求
func (s *Server) handleRequest(ctx context.Context, data []byte) Response {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return s.errorResponse(nil, ParseError, "解析请求失败", err)
	}

	// 验证 JSON-RPC 版本
	if req.JSONRPC != "2.0" {
		return s.errorResponse(req.ID, InvalidRequest, "无效的 JSON-RPC 版本", nil)
	}

	// 路由到相应的处理方法
	switch req.Method {
	case "initialize":
		return s.handleInitialize(ctx, req)
	case "initialized":
		// 客户端通知，不需要响应
		return Response{}
	case "tools/list":
		return s.handleListTools(ctx, req)
	case "tools/call":
		return s.handleCallTool(ctx, req)
	case "resources/list":
		return s.handleListResources(ctx, req)
	case "resources/read":
		return s.handleReadResource(ctx, req)
	case "prompts/list":
		return s.handleListPrompts(ctx, req)
	case "prompts/get":
		return s.handleGetPrompt(ctx, req)
	case "ping":
		return s.successResponse(req.ID, map[string]interface{}{})
	default:
		return s.errorResponse(req.ID, MethodNotFound, fmt.Sprintf("方法未找到: %s", req.Method), nil)
	}
}

// handleInitialize 处理初始化
func (s *Server) handleInitialize(ctx context.Context, req Request) Response {
	var params InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, InvalidParams, "无效的参数", err)
	}

	log.Printf("客户端初始化: %s v%s", params.ClientInfo.Name, params.ClientInfo.Version)

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo:      s.info,
	}

	return s.successResponse(req.ID, result)
}

// handleListTools 处理工具列表
func (s *Server) handleListTools(ctx context.Context, req Request) Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := ListToolsResult{
		Tools: s.tools,
	}

	return s.successResponse(req.ID, result)
}

// handleCallTool 处理调用工具
func (s *Server) handleCallTool(ctx context.Context, req Request) Response {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, InvalidParams, "无效的参数", err)
	}

	s.mu.RLock()
	handler, exists := s.toolHandlers[params.Name]
	s.mu.RUnlock()

	if !exists {
		return s.errorResponse(req.ID, InvalidParams, fmt.Sprintf("工具未找到: %s", params.Name), nil)
	}

	log.Printf("调用工具: %s", params.Name)

	result, err := handler(ctx, params.Arguments)
	if err != nil {
		return s.errorResponse(req.ID, InternalError, fmt.Sprintf("工具执行失败: %v", err), err)
	}

	return s.successResponse(req.ID, result)
}

// handleListResources 处理资源列表
func (s *Server) handleListResources(ctx context.Context, req Request) Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := ListResourcesResult{
		Resources: s.resources,
	}

	return s.successResponse(req.ID, result)
}

// handleReadResource 处理读取资源
func (s *Server) handleReadResource(ctx context.Context, req Request) Response {
	var params ReadResourceParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, InvalidParams, "无效的参数", err)
	}

	s.mu.RLock()
	handler, exists := s.resourceHandlers[params.URI]
	s.mu.RUnlock()

	if !exists {
		return s.errorResponse(req.ID, InvalidParams, fmt.Sprintf("资源未找到: %s", params.URI), nil)
	}

	log.Printf("读取资源: %s", params.URI)

	result, err := handler(ctx, params.URI)
	if err != nil {
		return s.errorResponse(req.ID, InternalError, fmt.Sprintf("资源读取失败: %v", err), err)
	}

	return s.successResponse(req.ID, result)
}

// handleListPrompts 处理提示列表
func (s *Server) handleListPrompts(ctx context.Context, req Request) Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := ListPromptsResult{
		Prompts: s.prompts,
	}

	return s.successResponse(req.ID, result)
}

// handleGetPrompt 处理获取提示
func (s *Server) handleGetPrompt(ctx context.Context, req Request) Response {
	var params GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, InvalidParams, "无效的参数", err)
	}

	s.mu.RLock()
	handler, exists := s.promptHandlers[params.Name]
	s.mu.RUnlock()

	if !exists {
		return s.errorResponse(req.ID, InvalidParams, fmt.Sprintf("提示未找到: %s", params.Name), nil)
	}

	log.Printf("获取提示: %s", params.Name)

	result, err := handler(ctx, params.Arguments)
	if err != nil {
		return s.errorResponse(req.ID, InternalError, fmt.Sprintf("提示生成失败: %v", err), err)
	}

	return s.successResponse(req.ID, result)
}

// successResponse 创建成功响应
func (s *Server) successResponse(id interface{}, result interface{}) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// errorResponse 创建错误响应
func (s *Server) errorResponse(id interface{}, code int, message string, err error) Response {
	errorData := &Error{
		Code:    code,
		Message: message,
	}

	if err != nil {
		errorData.Data = err.Error()
	}

	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   errorData,
	}
}
