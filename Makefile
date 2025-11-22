.PHONY: build run run-mcp clean docker docker-run docker-stop test

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
	@echo "  make build       - 构建应用"
	@echo "  make run         - 运行应用（HTTP 模式）"
	@echo "  make run-mcp     - 运行应用（MCP 模式）"
	@echo "  make clean       - 清理构建文件"
	@echo "  make docker      - 构建 Docker 镜像"
	@echo "  make docker-run  - 运行 Docker 容器"
	@echo "  make docker-stop - 停止 Docker 容器"
	@echo "  make test        - 运行测试"
	@echo "  make deps        - 下载依赖"
	@echo "  make fmt         - 格式化代码"
	@echo "  make lint        - 代码检查"
