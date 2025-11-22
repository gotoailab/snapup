# SnapUp MCP Server 使用指南

## 简介

SnapUp 现在支持 MCP (Model Context Protocol) 模式，可以作为大模型的工具服务，为 AI 助手提供网页截图能力。

## 什么是 MCP？

MCP (Model Context Protocol) 是 Anthropic 开发的一个开放标准协议，用于让大模型（如 Claude）能够安全地连接和使用外部工具、数据源和服务。通过 MCP，AI 助手可以：

- 调用外部工具和 API
- 访问实时数据
- 与各种服务集成
- 扩展其能力范围

## 运行模式

SnapUp 支持两种运行模式：

### 1. HTTP 模式（默认）
传统的 Web API 服务器，通过 HTTP 接口提供截图服务。

```bash
# 运行 HTTP 模式
./snapup -mode=http -port=8080

# 或使用 make
make run
```

### 2. MCP 模式
作为 MCP 服务器运行，通过标准输入/输出（stdio）与大模型通信。

```bash
# 运行 MCP 模式
./snapup -mode=mcp -output=./screenshots

# 或使用 make
make run-mcp
```

## MCP 配置

### 在 Claude Desktop 中配置

编辑 Claude Desktop 的配置文件：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

添加以下配置：

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

### 在其他 MCP 客户端中配置

任何支持 MCP 协议的客户端都可以使用 SnapUp。配置方式类似，需要提供：

- **command**: snapup 可执行文件的路径
- **args**: 命令行参数数组
  - `-mode=mcp`: 启用 MCP 模式
  - `-output=<path>`: 截图保存目录

## 可用工具

SnapUp MCP Server 提供以下工具：

### 1. take_screenshot

获取指定网站的屏幕截图。

**参数：**

- `url` (必需, string): 要截图的网站 URL
- `device` (可选, string): 设备类型
  - `desktop`: 桌面 (1920x1080) - 默认
  - `laptop`: 笔记本 (1440x900)
  - `tablet`: 平板 (768x1024)
  - `mobile`: 手机 (375x812)
- `style` (可选, string): 截图样式
  - `none`: 无样式 - 默认
  - `glass`: 玻璃风格
  - `device`: 设备边框
  - `floating`: 浮动阴影
- `full_page` (可选, boolean): 是否全页截图，默认 false
- `delay` (可选, integer): 截图前延迟（毫秒），默认 1000，范围 0-30000
- `quality` (可选, integer): 图片质量，默认 90，范围 1-100
- `background` (可选, string): 背景颜色，默认 "#f0f2f5"

**返回：**
- 文本描述（包含截图信息）
- Base64 编码的 PNG 图片

**示例对话：**

```
用户: 帮我截取 https://www.example.com 在手机上的样子

AI: 好的，我来帮你截取该网站在手机设备上的截图...
[调用 take_screenshot 工具]
[返回手机尺寸的网页截图]
```

### 2. get_devices_info

获取所有支持的设备类型及其屏幕尺寸信息。

**参数：** 无

**返回：**
设备类型列表，包含名称、尺寸等信息。

### 3. get_styles_info

获取所有支持的截图样式及其描述。

**参数：** 无

**返回：**
样式列表，包含名称和描述。

## 使用示例

### 示例 1: 基本截图

```
用户: 截取 https://github.com 的桌面版截图

AI: [调用 take_screenshot]
参数:
- url: "https://github.com"
- device: "desktop"
- style: "none"
```

### 示例 2: 移动端全页截图

```
用户: 给我看看 https://www.apple.com 在 iPhone 上的完整页面

AI: [调用 take_screenshot]
参数:
- url: "https://www.apple.com"
- device: "mobile"
- full_page: true
- style: "device"
```

### 示例 3: 多设备对比

```
用户: 帮我对比 https://www.example.com 在不同设备上的显示效果

AI: 好的，我将为你截取桌面、平板和手机三种设备的截图...
[依次调用 take_screenshot 三次，使用不同的 device 参数]
```

### 示例 4: 带样式的截图

```
用户: 截取 https://www.google.com 的玻璃风格截图

AI: [调用 take_screenshot]
参数:
- url: "https://www.google.com"
- device: "desktop"
- style: "glass"
```

## 高级配置

### 自定义截图目录

```bash
./snapup -mode=mcp -output=/custom/path/screenshots
```

### 性能优化

对于需要频繁截图的场景，建议：

1. 使用 SSD 存储截图文件
2. 适当增加系统内存
3. 确保网络连接稳定

## 故障排除

### 问题 1: MCP Server 无法启动

**可能原因：**
- 截图目录没有写入权限
- 端口被占用（HTTP 模式）
- Chrome/Chromium 未安装

**解决方法：**
```bash
# 检查截图目录权限
mkdir -p ./screenshots
chmod 755 ./screenshots

# 安装 Chrome/Chromium (Linux)
sudo apt-get install chromium-browser
```

### 问题 2: 截图失败

**可能原因：**
- URL 无法访问
- 页面加载超时
- 网络连接问题

**解决方法：**
- 增加 delay 参数值
- 检查 URL 是否正确
- 确认网络连接

### 问题 3: 图片质量不佳

**解决方法：**
- 增加 quality 参数值（最大 100）
- 选择合适的设备分辨率
- 使用 full_page: false 只截取可见区域

## 开发集成

### 直接调用 JSON-RPC

MCP Server 使用 JSON-RPC 2.0 协议，可以通过标准输入/输出进行通信：

```json
// 请求
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "take_screenshot",
    "arguments": {
      "url": "https://www.example.com",
      "device": "desktop",
      "style": "none"
    }
  }
}

// 响应
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "截图成功！..."
      },
      {
        "type": "image",
        "data": "base64_encoded_image_data...",
        "mimeType": "image/png"
      }
    ],
    "isError": false
  }
}
```

## 安全注意事项

1. **访问控制**: 确保只有授权的客户端可以连接到 MCP Server
2. **资源限制**: 截图操作可能消耗大量系统资源，建议设置合理的并发限制
3. **存储管理**: 定期清理旧的截图文件，避免磁盘空间耗尽
4. **URL 验证**: 谨慎处理用户提供的 URL，防止访问恶意网站

## 更多信息

- [SnapUp GitHub 仓库](https://github.com/gotoailab/snapup)
- [MCP 协议规范](https://modelcontextprotocol.io)
- [Anthropic Claude Desktop](https://claude.ai/desktop)

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

