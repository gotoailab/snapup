# SnapUp MCP 快速开始

本指南帮助你在 5 分钟内将 SnapUp 集成到 Claude Desktop 或其他 MCP 客户端中。

## 第 1 步：构建 SnapUp

```bash
cd /path/to/snapup
go build -o snapup ./cmd/snapup
```

或使用 Make：

```bash
make build
```

## 第 2 步：测试 MCP 模式

```bash
./snapup -mode=mcp -output=./screenshots
```

如果看到 "MCP Server 启动，使用 stdio 传输" 消息，说明 MCP 模式工作正常。按 Ctrl+C 退出。

## 第 3 步：配置 Claude Desktop

### macOS

编辑配置文件：

```bash
nano ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

### Windows

编辑配置文件：

```
%APPDATA%\Claude\claude_desktop_config.json
```

### Linux

编辑配置文件：

```bash
nano ~/.config/Claude/claude_desktop_config.json
```

### 配置内容

```json
{
  "mcpServers": {
    "snapup": {
      "command": "/absolute/path/to/snapup",
      "args": ["-mode=mcp", "-output=/absolute/path/to/screenshots"]
    }
  }
}
```

**重要提示：**
- 必须使用**绝对路径**
- 确保 screenshots 目录已存在且有写入权限
- 如果已有其他 MCP 服务器配置，将 snapup 配置添加到 mcpServers 对象中

### 示例（macOS）

```json
{
  "mcpServers": {
    "snapup": {
      "command": "/Users/username/projects/snapup/snapup",
      "args": ["-mode=mcp", "-output=/Users/username/projects/snapup/screenshots"]
    }
  }
}
```

## 第 4 步：重启 Claude Desktop

1. 完全退出 Claude Desktop（不是最小化）
2. 重新启动 Claude Desktop
3. 检查是否成功连接（可以在 Claude 的设置或日志中查看）

## 第 5 步：测试

在 Claude Desktop 中尝试以下对话：

```
你好！请帮我截取 https://www.google.com 的桌面版截图。
```

Claude 应该会调用 `take_screenshot` 工具，并返回截图结果。

## 常用命令示例

### 基本截图

```
截取 https://github.com 的截图
```

### 移动设备截图

```
帮我看看 https://www.apple.com 在 iPhone 上的样子
```

### 多设备对比

```
对比 https://example.com 在桌面、平板和手机上的显示效果
```

### 全页截图

```
截取 https://news.ycombinator.com 的完整页面
```

### 带样式的截图

```
用玻璃风格截取 https://www.stripe.com
```

## 支持的设备类型

- `desktop` - 桌面 (1920x1080)
- `laptop` - 笔记本 (1440x900)
- `tablet` - 平板 (768x1024)
- `mobile` - 手机 (375x812)

## 支持的样式

- `none` - 无样式（默认）
- `glass` - 玻璃风格
- `device` - 设备边框
- `floating` - 浮动阴影

## 故障排除

### 问题：Claude 找不到 snapup 工具

**解决方法：**
1. 检查配置文件中的路径是否正确（使用绝对路径）
2. 确认 snapup 可执行文件有执行权限：`chmod +x /path/to/snapup`
3. 重启 Claude Desktop

### 问题：截图失败

**解决方法：**
1. 确保 screenshots 目录存在且有写入权限
2. 确保系统已安装 Chrome 或 Chromium
3. 检查网络连接

### 问题：查看详细日志

在配置文件中添加环境变量：

```json
{
  "mcpServers": {
    "snapup": {
      "command": "/path/to/snapup",
      "args": ["-mode=mcp", "-output=/path/to/screenshots"],
      "env": {
        "DEBUG": "true"
      }
    }
  }
}
```

## 下一步

- 查看 [MCP_USAGE.md](./MCP_USAGE.md) 了解完整功能
- 阅读 [README.md](./README.md) 了解项目详情
- 参考 [examples/](./examples/) 目录中的示例

## 需要帮助？

- 提交 Issue: [GitHub Issues](https://github.com/gotoailab/snapup/issues)
- 查看文档: [MCP Protocol](https://modelcontextprotocol.io)

祝你使用愉快！🚀

