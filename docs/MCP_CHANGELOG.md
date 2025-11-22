# MCP 功能变更日志

## 版本 1.1.0 - 添加 MCP Server 支持

**发布日期**: 2025-11-22

### 🎉 新增功能

#### MCP (Model Context Protocol) 支持

SnapUp 现在可以作为 MCP 服务器运行，为大型语言模型提供网页截图工具能力。

### 📝 新增文件

#### 核心实现

1. **internal/mcp/types.go** - MCP 协议类型定义
   - JSON-RPC 2.0 请求/响应类型
   - MCP 协议数据结构
   - 工具、资源、提示相关类型

2. **internal/mcp/server.go** - MCP 服务器核心实现
   - JSON-RPC 2.0 协议处理
   - stdio 通信支持
   - 工具/资源/提示注册和调度

3. **internal/mcp/tools.go** - 截图工具封装
   - `take_screenshot` - 网页截图工具
   - `get_devices_info` - 设备信息查询工具
   - `get_styles_info` - 样式信息查询工具

#### 文档

4. **MCP_USAGE.md** - MCP 详细使用指南
   - MCP 概念介绍
   - 配置说明
   - 工具参考文档
   - 使用示例
   - 故障排除

5. **MCP_QUICKSTART.md** - MCP 快速开始指南
   - 5 分钟快速上手
   - Claude Desktop 配置步骤
   - 常见问题解决

6. **MCP_CHANGELOG.md** - 本文件，MCP 功能变更日志

#### 示例

7. **examples/mcp_config_example.json** - MCP 配置示例
8. **examples/mcp_test.sh** - MCP 测试脚本

### 🔧 修改的文件

#### 主程序

1. **cmd/snapup/main.go**
   - 添加 `-mode` 参数（http/mcp）
   - 添加 `runMCPServer()` 函数
   - 添加信号处理支持

#### 构建配置

2. **Makefile**
   - 添加 `run-mcp` 命令
   - 更新 `help` 命令说明
   - 更新 `.PHONY` 目标

#### 文档

3. **README.md**
   - 在特性列表中添加 MCP 支持
   - 添加 MCP 模式快速开始章节
   - 更新项目结构说明
   - 更新构建命令说明

### 🛠️ 技术细节

#### 运行模式

SnapUp 现在支持两种运行模式：

1. **HTTP 模式**（默认）
   ```bash
   ./snapup -mode=http -port=8080
   ```
   传统的 Web API 服务器

2. **MCP 模式**（新增）
   ```bash
   ./snapup -mode=mcp -output=./screenshots
   ```
   通过 stdio 提供 MCP 服务

#### MCP 工具

##### take_screenshot

最核心的工具，提供网页截图功能。

**输入参数**:
- `url` (必需): 目标网址
- `device`: 设备类型（desktop/laptop/tablet/mobile）
- `style`: 样式效果（none/glass/device/floating）
- `full_page`: 是否全页截图
- `delay`: 延迟时间（毫秒）
- `quality`: 图片质量（1-100）
- `background`: 背景颜色

**输出**:
- 文本描述（截图信息）
- Base64 编码的 PNG 图片

##### get_devices_info

查询支持的设备类型信息。

##### get_styles_info

查询支持的样式效果信息。

#### 协议实现

- JSON-RPC 2.0 over stdio
- MCP Protocol Version: 2024-11-05
- 支持的功能：
  - ✅ Tools (工具)
  - ✅ Resources (资源) - 框架支持
  - ✅ Prompts (提示) - 框架支持

### 📦 依赖变化

无新增外部依赖。所有 MCP 功能使用 Go 标准库实现。

### 🔄 向后兼容性

完全向后兼容。现有的 HTTP API 和功能不受影响。

- 默认运行模式仍为 HTTP
- 所有现有命令和 API 保持不变
- Docker 部署方式不变

### 🚀 使用方式

#### 在 Claude Desktop 中使用

1. 构建 SnapUp
2. 编辑 Claude Desktop 配置文件
3. 添加 MCP 服务器配置
4. 重启 Claude Desktop
5. 在对话中使用截图功能

详见 [MCP_QUICKSTART.md](./MCP_QUICKSTART.md)

#### 编程方式集成

任何支持 MCP 协议的客户端都可以使用 SnapUp 的截图功能。

### 🐛 已知问题

无

### 📋 未来计划

- [ ] 添加更多截图选项（例如：自定义视口大小）
- [ ] 支持批量截图
- [ ] 添加截图对比功能
- [ ] 支持更多输出格式（JPEG, WebP）
- [ ] 添加截图缓存机制

### 🙏 致谢

感谢 Anthropic 开发的 MCP 协议，让 AI 工具集成变得简单而强大。

---

## 升级指南

### 从 1.0.x 升级到 1.1.0

1. **拉取最新代码**
   ```bash
   git pull origin main
   ```

2. **重新构建**
   ```bash
   make build
   # 或
   go build -o snapup ./cmd/snapup
   ```

3. **（可选）配置 MCP 模式**
   
   如果要使用 MCP 功能，按照 [MCP_QUICKSTART.md](./MCP_QUICKSTART.md) 进行配置。

4. **运行**
   ```bash
   # HTTP 模式（与之前完全相同）
   ./snapup -port=8080
   
   # 或 MCP 模式（新功能）
   ./snapup -mode=mcp
   ```

### 配置迁移

无需任何配置迁移。现有配置和部署方式继续有效。

---

如有问题或建议，请提交 Issue 到 GitHub 仓库。

