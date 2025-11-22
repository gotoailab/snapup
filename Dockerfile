# 构建阶段
FROM golang:1.21-alpine AS builder

# 构建参数
ARG GOPROXY=""

# 设置工作目录
WORKDIR /app

# 如果提供了 GOPROXY，则设置 Go 代理
ENV GOPROXY=${GOPROXY}
ENV GO111MODULE=on

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o snapup ./cmd/snapup

# 运行阶段
FROM debian:bookworm-slim

# 仅安装运行时必需的基础依赖
# 注意：不需要安装 Chrome，因为使用独立的 Chrome 容器
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/snapup .

# 创建截图输出目录
RUN mkdir -p /app/screenshots

# 环境变量配置
ENV RUN_MODE=http
ENV SERVER_PORT=8080
ENV OUTPUT_DIR=/app/screenshots

# 暴露端口
EXPOSE 8080

# 运行应用（使用 shell 形式以支持环境变量）
CMD ./snapup -mode=${RUN_MODE} -port=${SERVER_PORT} -output=${OUTPUT_DIR}
