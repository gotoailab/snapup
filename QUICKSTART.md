# 快速开始指南

5 分钟内启动并运行 SnapUp 截图服务。

## 方式一：使用 Docker（推荐）

### 步骤 1: 启动服务

```bash
# 确保你在项目目录下
cd snapup

# 构建并启动服务
docker-compose up -d
```

### 步骤 2: 验证服务

```bash
# 检查服务状态
curl http://localhost:8080/api/health

# 应该看到: {"status":"ok","service":"snapup"}
```

### 步骤 3: 访问 Web 界面

打开浏览器访问: http://localhost:8080

### 步骤 4: 生成第一张截图

在 Web 界面中：
1. 输入网址: `https://example.com`
2. 选择设备: 桌面
3. 选择样式: 玻璃风格
4. 点击"生成截图"

完成！你应该能看到生成的截图。

## 方式二：本地运行

### 前置要求

- Go 1.21+
- Google Chrome

### 步骤 1: 安装依赖

```bash
# 下载 Go 依赖
go mod download
```

### 步骤 2: 构建项目

```bash
# 使用 Makefile
make build

# 或者手动构建
go build -o snapup ./cmd/snapup
```

### 步骤 3: 运行服务

```bash
# 使用 Makefile
make run

# 或者直接运行
./snapup -port 8080
```

### 步骤 4: 访问服务

打开浏览器访问: http://localhost:8080

## API 快速测试

### 使用 curl

```bash
# 基础截图
curl -X POST http://localhost:8080/api/screenshot \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "device": "desktop",
    "style": "glass"
  }'
```

### 使用测试脚本

```bash
cd examples
./api_test.sh
```

## 常用命令

### Docker 方式

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重新构建
docker-compose up -d --build
```

### 本地方式

```bash
# 构建
make build

# 运行
make run

# 清理
make clean

# 格式化代码
make fmt
```

## 配置选项

### 命令行参数

```bash
./snapup -port 8080 -output ./screenshots
```

- `-port`: 服务器端口（默认：8080）
- `-output`: 截图输出目录（默认：./screenshots）

## 使用 Web 界面

### 基本流程

1. **输入网址** - 在输入框输入要截图的 URL
2. **选择设备** - 桌面、笔记本、平板或手机
3. **选择样式** - 无样式、玻璃风格、设备边框或浮动阴影
4. **配置选项**（可选）
   - 全页截图
   - 延迟时间
   - 背景颜色
   - 图片质量
5. **生成截图** - 点击按钮生成
6. **预览和下载** - 查看结果并下载图片

## 故障排除

### 端口被占用

```bash
# 更改端口
./snapup -port 8090

# 或修改 docker-compose.yml
ports:
  - "8090:8080"
```

### Chrome 未找到

确保 Chrome 已安装：

```bash
# Ubuntu/Debian
sudo apt-get install google-chrome-stable

# 或使用 Docker（推荐）
docker-compose up -d
```

### 权限错误

```bash
# 创建输出目录并设置权限
mkdir -p screenshots
chmod 755 screenshots
```

## 下一步

- 阅读 [README.md](README.md) 了解详细功能
- 查看 [examples/](examples/) 目录查看更多示例
- 查看 API 文档了解所有可用选项

享受使用 SnapUp！
