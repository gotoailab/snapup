.PHONY: build run run-mcp clean docker docker-run docker-stop docker-cn docker-run-cn docker-stop-cn docker-mcp-start docker-mcp-stop docker-mcp-restart docker-mcp-logs docker-mcp-status docker-mcp-build docker-mcp-start-cn test deps fmt lint help

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
