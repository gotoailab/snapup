# SnapUp - 专业网页截图生成器

SnapUp 是一个基于 Go 和 ChromeDP 开发的高性能网页截图服务，支持多种设备类型和样式效果，为您的网页生成精美的截图 Mockup。

## 特性

- 🚀 **高性能**: 使用 Go 和 ChromeDP 提供快速的截图服务
- 📱 **多设备支持**: 支持桌面、笔记本、平板和手机等多种设备尺寸
- 🎨 **多样式效果**: 提供玻璃风格、设备边框、浮动阴影等多种样式
- 🔧 **高度可配置**: 支持自定义延迟、背景颜色、图片质量等参数
- 📄 **全页截图**: 支持捕获完整网页内容
- 🐳 **Docker 支持**: 提供完整的 Docker 部署方案
- 💻 **现代化界面**: 使用 Vue 3 和 Tailwind CSS 构建的美观界面
- 🤖 **MCP 支持**: 支持 Model Context Protocol，可作为大模型工具使用

## 快速开始

### 使用 Docker（推荐）

这是最简单的部署方式，Docker Compose 会自动启动 Chrome headless 容器和应用服务。

```bash
# 克隆项目
git clone <repository-url>
cd snapup

# 构建并运行（会自动启动 Chrome 和 SnapUp 服务）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 访问服务
# 浏览器打开 http://localhost:8080

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

**架构说明**：
- `chrome` 服务：运行 Chrome headless 浏览器，监听 9222 端口
- `snapup` 服务：运行主应用，通过 WebSocket 连接到 Chrome 容器
- 自动健康检查，确保 Chrome 启动后才启动主服务

### 本地运行

#### 前置要求

- Go 1.21 或更高版本
- Google Chrome 或 Chromium 浏览器

#### 安装步骤

1. **克隆项目**

```bash
git clone <repository-url>
cd snapup
```

2. **安装依赖**

```bash
go mod download
```

3. **构建项目**

```bash
make build
# 或者
go build -o snapup ./cmd/snapup
```

4. **运行服务**

```bash
./snapup -port 8080 -output ./screenshots
# 或者使用 make
make run
```

5. **访问服务**

浏览器打开 `http://localhost:8080`

### MCP 模式（AI 工具集成）

SnapUp 支持作为 MCP (Model Context Protocol) 服务器运行，可以被 Claude Desktop 等 AI 助手调用，为大模型提供网页截图能力。

```bash
# 运行 MCP 模式
./snapup -mode=mcp -output ./screenshots

# 或使用 make
make run-mcp
```

**配置 Claude Desktop:**

编辑配置文件（macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`）：

```json
{
  "mcpServers": {
    "snapup": {
      "command": "/path/to/snapup",
      "args": ["-mode=mcp", "-output=/path/to/screenshots"]
    }
  }
}
```

配置完成后，重启 Claude Desktop，你就可以在对话中要求 Claude 帮你截取网页了！

**详细使用指南**: 查看 [MCP_USAGE.md](./MCP_USAGE.md) 了解完整的 MCP 功能和使用方法。

## 配置

### 环境变量

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `PORT` | 服务监听端口 | `8080` | `8080` |
| `CHROME_WS_URL` | Chrome WebSocket URL（用于 Docker 部署） | 空（使用本地 Chrome） | `ws://chrome:9222` |

**Docker 部署**：`CHROME_WS_URL` 会自动配置为 `ws://chrome:9222`，连接到 Chrome 容器

**本地部署**：不设置 `CHROME_WS_URL`，程序会自动启动本地 Chrome 实例

### Docker Compose 配置

默认的 `docker-compose.yml` 包含两个服务：

1. **chrome** - Chrome headless 容器
   - 端口：9222（WebSocket 调试端口）
   - 内存：2GB shared memory
   - 健康检查：自动验证 Chrome 是否就绪

2. **snapup** - 主应用服务
   - 端口：8080（HTTP 服务）
   - 依赖：等待 Chrome 容器健康后启动
   - 卷挂载：`./screenshots` 目录用于存储截图

## 使用说明

### Web 界面

1. 在输入框中输入要截图的网址
2. 选择设备类型（桌面/笔记本/平板/手机）
3. 选择样式效果（无样式/玻璃风格/设备边框/浮动阴影）
4. 配置高级选项（可选）
   - 全页截图
   - 延迟时间
   - 背景颜色
   - 图片质量
5. 点击"生成截图"按钮
6. 等待生成完成后可预览和下载

### API 接口

#### 生成截图

**请求**

```http
POST /api/screenshot
Content-Type: application/json

{
  "url": "https://example.com",
  "device": "desktop",
  "style": "glass",
  "delay": 1000,
  "full_page": false,
  "quality": 90,
  "background": "#f0f2f5"
}
```

**参数说明**

| 参数 | 类型 | 说明 | 可选值 |
|------|------|------|--------|
| url | string | 要截图的网址 | 任意有效 URL |
| device | string | 设备类型 | desktop, laptop, tablet, mobile |
| style | string | 样式效果 | none, glass, device, floating |
| delay | int | 延迟时间(毫秒) | 0-10000 |
| full_page | bool | 是否全页截图 | true, false |
| quality | int | 图片质量 | 1-100 |
| background | string | 背景颜色 | 十六进制颜色值 |

**响应**

```json
{
  "success": true,
  "message": "截图成功",
  "image_url": "/screenshots/screenshot_desktop_glass_xxx.png",
  "filename": "screenshot_desktop_glass_xxx.png"
}
```

## 项目结构

```
snapup/
├── cmd/
│   └── snapup/          # 主程序入口
│       └── main.go
├── internal/
│   ├── mcp/             # MCP 服务器实现
│   │   ├── types.go     # MCP 协议类型
│   │   ├── server.go    # MCP 服务器
│   │   └── tools.go     # 截图工具封装
│   ├── models/          # 数据模型
│   │   └── types.go
│   ├── screenshot/      # 截图核心功能
│   │   ├── capture.go   # 截图捕获
│   │   ├── processor.go # 图片处理
│   │   └── service.go   # 截图服务
│   └── server/          # Web 服务器
│       ├── handler.go   # HTTP 处理器
│       ├── middleware.go # 中间件
│       ├── server.go    # 服务器
│       └── static/      # 静态文件
│           └── index.html
├── screenshots/         # 截图输出目录
├── Dockerfile          # Docker 构建文件
├── docker-compose.yml  # Docker Compose 配置
├── Makefile           # Make 命令
├── go.mod             # Go 模块定义
├── README.md          # 项目说明
└── MCP_USAGE.md       # MCP 使用指南
```

## 开发指南

### 构建命令

```bash
# 构建
make build

# 运行（HTTP 模式）
make run

# 运行（MCP 模式）
make run-mcp

# 清理
make clean

# 格式化代码
make fmt

# 代码检查
make lint

# 运行测试
make test
```

### Docker 命令

```bash
# 构建镜像
make docker

# 启动容器
make docker-run

# 停止容器
make docker-stop
```

## 技术栈

- **后端**: Go 1.21+
- **截图引擎**: ChromeDP
- **前端框架**: Vue 3
- **样式框架**: Tailwind CSS
- **容器化**: Docker & Docker Compose

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
