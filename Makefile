.PHONY: build run run-mcp clean docker docker-run docker-stop docker-cn docker-run-cn docker-stop-cn docker-mcp-start docker-mcp-stop docker-mcp-restart docker-mcp-logs docker-mcp-status docker-mcp-build docker-mcp-start-cn test deps fmt lint help tunnel-setup tunnel-test tunnel-frp tunnel-cloudflare tunnel-nginx tunnel-stop

# 构建应用
build:
	go build -o snapup ./cmd/snapup

# 运行应用（HTTP 模式）
run:
	go run ./cmd/snapup/main.go -mode=http

# 运行应用（MCP 模式）
run-mcp:
	go run ./cmd/snapup/main.go -mode=mcp

# 清理
clean:
	rm -f snapup
	rm -rf screenshots/*.png

# 构建 Docker 镜像
docker:
	docker build -t snapup:latest .

# 运行 Docker 容器
docker-run:
	docker-compose up -d

# 停止 Docker 容器
docker-stop:
	docker-compose down

# 构建 Docker 镜像（中国版）
docker-cn:
	docker build -f Dockerfile.cn -t snapup:latest .

# 运行 Docker 容器（中国版）
docker-run-cn:
	docker-compose -f docker-compose.cn.yml up -d

# 停止 Docker 容器（中国版）
docker-stop-cn:
	docker-compose -f docker-compose.cn.yml down

# 启动 Docker MCP 服务
docker-mcp-start:
	@./scripts/mcp-docker.sh start

# 停止 Docker MCP 服务
docker-mcp-stop:
	@./scripts/mcp-docker.sh stop

# 查看 Docker MCP 服务日志
docker-mcp-logs:
	@./scripts/mcp-docker.sh logs -f

# 查看 Docker MCP 服务状态
docker-mcp-status:
	@./scripts/mcp-docker.sh status

# 启动 Docker MCP 服务（中国版）
docker-mcp-start-cn:
	@./scripts/mcp-docker.sh start --cn

# 重启 Docker MCP 服务
docker-mcp-restart:
	@./scripts/mcp-docker.sh restart

# 重新构建 Docker MCP 服务
docker-mcp-build:
	@./scripts/mcp-docker.sh build

# 运行测试
test:
	go test -v ./...

# 下载依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	go vet ./...

# 内网穿透配置向导
tunnel-setup:
	@./setup-tunnel.sh

# 测试内网穿透连接
tunnel-test:
	@if [ -z "$(URL)" ]; then \
		echo "用法: make tunnel-test URL=http://your-server-ip:9222"; \
		echo "或: make tunnel-test URL=https://chrome.your-domain.com"; \
		exit 1; \
	fi
	@./test-tunnel.sh $(URL)

# 启动 frp 内网穿透
tunnel-frp:
	docker-compose -f docker-compose.tunnel.yml up -d
	@echo "frp 内网穿透已启动"
	@echo "查看日志: make tunnel-logs"

# 启动 Cloudflare Tunnel
tunnel-cloudflare:
	docker-compose -f docker-compose.cloudflare.yml up -d
	@echo "Cloudflare Tunnel 已启动"
	@echo "查看日志: make tunnel-logs"

# 启动 Nginx + 内网穿透
tunnel-nginx:
	docker-compose -f docker-compose.nginx-tunnel.yml up -d
	@echo "Nginx + 内网穿透已启动"
	@echo "查看日志: make tunnel-logs"

# 查看内网穿透日志
tunnel-logs:
	@if docker ps | grep -q snapup-frpc; then \
		docker-compose -f docker-compose.tunnel.yml logs -f frpc; \
	elif docker ps | grep -q snapup-cloudflared; then \
		docker-compose -f docker-compose.cloudflare.yml logs -f cloudflared; \
	elif docker ps | grep -q snapup-nginx-ws; then \
		docker-compose -f docker-compose.nginx-tunnel.yml logs -f; \
	else \
		echo "没有运行中的内网穿透服务"; \
	fi

# 停止内网穿透服务
tunnel-stop:
	@echo "停止所有内网穿透服务..."
	@docker-compose -f docker-compose.tunnel.yml down 2>/dev/null || true
	@docker-compose -f docker-compose.cloudflare.yml down 2>/dev/null || true
	@docker-compose -f docker-compose.nginx-tunnel.yml down 2>/dev/null || true
	@echo "内网穿透服务已停止"

# 帮助
help:
	@echo "可用命令:"
	@echo ""
	@echo "基本命令:"
	@echo "  make build              - 构建应用"
	@echo "  make run                - 运行应用（HTTP 模式）"
	@echo "  make run-mcp            - 运行应用（MCP 模式）"
	@echo "  make clean              - 清理构建文件"
	@echo ""
	@echo "Docker 命令:"
	@echo "  make docker             - 构建 Docker 镜像"
	@echo "  make docker-run         - 运行 Docker 容器（HTTP 模式）"
	@echo "  make docker-stop        - 停止 Docker 容器"
	@echo "  make docker-cn          - 构建 Docker 镜像（中国版）"
	@echo "  make docker-run-cn      - 运行 Docker 容器（中国版）"
	@echo "  make docker-stop-cn     - 停止 Docker 容器（中国版）"
	@echo ""
	@echo "Docker MCP 命令:"
	@echo "  make docker-mcp-start   - 启动 Docker MCP 服务"
	@echo "  make docker-mcp-stop    - 停止 Docker MCP 服务"
	@echo "  make docker-mcp-restart - 重启 Docker MCP 服务"
	@echo "  make docker-mcp-logs    - 查看 Docker MCP 服务日志"
	@echo "  make docker-mcp-status  - 查看 Docker MCP 服务状态"
	@echo "  make docker-mcp-build   - 重新构建 Docker MCP 服务"
	@echo "  make docker-mcp-start-cn- 启动 Docker MCP 服务（中国版）"
	@echo ""
	@echo "开发命令:"
	@echo "  make test               - 运行测试"
	@echo "  make deps               - 下载依赖"
	@echo "  make fmt                - 格式化代码"
	@echo "  make lint               - 代码检查"
	@echo ""
	@echo "内网穿透命令:"
	@echo "  make tunnel-setup       - 运行配置向导"
	@echo "  make tunnel-frp         - 启动 frp 内网穿透"
	@echo "  make tunnel-cloudflare  - 启动 Cloudflare Tunnel"
	@echo "  make tunnel-nginx       - 启动 Nginx + 内网穿透"
	@echo "  make tunnel-logs        - 查看内网穿透日志"
	@echo "  make tunnel-stop        - 停止所有内网穿透服务"
	@echo "  make tunnel-test URL=<地址> - 测试内网穿透连接"
	@echo ""
	@echo "示例:"
	@echo "  make tunnel-test URL=http://your-server-ip:9222"
	@echo "  make tunnel-test URL=https://chrome.your-domain.com"
