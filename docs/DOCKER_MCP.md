# Docker 运行 MCP Server 指南

## 简介

SnapUp 的 MCP Server 现在可以通过 Docker 容器运行，使用 Docker Compose 进行编排。

## 快速开始

### 1. 启动 MCP Server（标准版）

```bash
# 构建并启动 MCP 服务（包含 Chrome 容器）
docker-compose --profile mcp up -d

# 查看 MCP 服务日志
docker-compose logs -f snapup-mcp
```

### 2. 启动 MCP Server（中国版 - 使用国内镜像源）

如果你在中国大陆，建议使用中国版配置文件，构建速度更快：

```bash
# 构建并启动 MCP 服务（使用国内镜像源）
docker-compose -f docker-compose.cn.yml --profile mcp up -d

# 查看 MCP 服务日志
docker-compose -f docker-compose.cn.yml logs -f snapup-mcp
```

## 配置说明

### Docker Compose Profile

MCP Server 使用了 Docker Compose 的 `profile` 特性，默认不会自动启动。你需要显式指定 `--profile mcp` 来启动它。

### 服务配置

MCP Server 容器的主要配置：

- **容器名**: `snapup-mcp`
- **重启策略**: `no`（不自动重启）
- **依赖服务**: `chrome`（无头浏览器）
- **输出目录**: `./screenshots`（挂载到主机）
- **stdin/tty**: 启用（MCP 协议需要）

### 环境变量

- `RUN_MODE=mcp`: 运行模式设置为 MCP
- `OUTPUT_DIR=/app/screenshots`: 截图保存目录
- `CHROME_WS_URL=ws://chrome:9222`: Chrome DevTools WebSocket URL

## 使用场景

### 场景 1: 同时运行 HTTP 和 MCP 模式

```bash
# 启动 HTTP 服务（默认）
docker-compose up -d

# 同时启动 MCP 服务
docker-compose --profile mcp up -d snapup-mcp

# 此时两个服务都在运行
docker-compose ps
```

### 场景 2: 仅运行 MCP 模式

```bash
# 只启动 Chrome 和 MCP 服务
docker-compose up -d chrome
docker-compose --profile mcp up -d snapup-mcp
```

### 场景 3: 临时测试 MCP 服务

```bash
# 以非后台模式运行，方便查看日志
docker-compose --profile mcp up snapup-mcp

# Ctrl+C 停止
```

## 连接 MCP Server

### 从宿主机连接

虽然 MCP Server 在容器中运行，但你仍然可以从宿主机通过 stdin/stdout 与它交互。

### 在 Claude Desktop 中配置

由于 MCP Server 在容器中运行，你需要使用 Docker 命令来启动它：

```json
{
  "mcpServers": {
    "snapup": {
      "command": "docker",
      "args": [
        "exec",
        "-i",
        "snapup-mcp",
        "/app/snapup",
        "-mode=mcp",
        "-output=/app/screenshots"
      ]
    }
  }
}
```

**注意**: 
- 确保 `snapup-mcp` 容器已经启动
- 使用 `docker exec -i` 来与容器交互
- `-i` 参数保持 stdin 打开

### 使用 Docker Compose 命令

更推荐使用 docker-compose 命令：

```json
{
  "mcpServers": {
    "snapup": {
      "command": "docker-compose",
      "args": [
        "exec",
        "-T",
        "snapup-mcp",
        "/app/snapup",
        "-mode=mcp",
        "-output=/app/screenshots"
      ],
      "cwd": "/path/to/snapup"
    }
  }
}
```

**参数说明**:
- `exec`: 在运行中的容器内执行命令
- `-T`: 禁用伪终端分配（stdin/stdout 模式需要）
- `cwd`: 工作目录，指向 docker-compose.yml 所在目录

## 常见命令

### 查看服务状态

```bash
# 查看所有服务
docker-compose ps

# 查看 MCP 服务
docker-compose ps snapup-mcp
```

### 查看日志

```bash
# 查看 MCP 服务日志
docker-compose logs snapup-mcp

# 实时跟踪日志
docker-compose logs -f snapup-mcp

# 查看最近 100 行日志
docker-compose logs --tail=100 snapup-mcp
```

### 重启服务

```bash
# 重启 MCP 服务
docker-compose restart snapup-mcp

# 停止后重新启动
docker-compose stop snapup-mcp
docker-compose --profile mcp up -d snapup-mcp
```

### 停止服务

```bash
# 停止 MCP 服务
docker-compose stop snapup-mcp

# 停止所有服务
docker-compose down
```

### 重新构建

```bash
# 重新构建 MCP 服务镜像
docker-compose build snapup-mcp

# 强制重新构建（不使用缓存）
docker-compose build --no-cache snapup-mcp

# 重新构建并启动
docker-compose --profile mcp up -d --build snapup-mcp
```

## 自定义构建参数

### 使用国内 Go 代理

Dockerfile 现在支持 `GOPROXY` 构建参数：

```bash
# 使用国内 Go 代理构建
docker-compose build --build-arg GOPROXY=https://goproxy.cn,direct snapup-mcp

# 或者编辑 docker-compose.yml，添加 build args
```

### 修改 docker-compose.yml

```yaml
snapup-mcp:
  build:
    context: .
    dockerfile: Dockerfile
    args:
      - GOPROXY=https://goproxy.cn,direct
  # ... 其他配置
```

## 故障排除

### 问题 1: MCP 服务无法启动

**症状**: `docker-compose up --profile mcp` 后 MCP 容器启动失败

**解决方案**:
```bash
# 检查 Chrome 容器是否运行
docker-compose ps chrome

# 如果 Chrome 未运行，先启动它
docker-compose up -d chrome

# 查看详细日志
docker-compose logs snapup-mcp
```

### 问题 2: 无法连接到 Chrome

**症状**: 日志显示无法连接到 WebSocket

**解决方案**:
```bash
# 检查 Chrome 容器健康状态
docker-compose ps

# 重启 Chrome 容器
docker-compose restart chrome

# 等待 10 秒后重启 MCP
sleep 10
docker-compose restart snapup-mcp
```

### 问题 3: 截图文件无法保存

**症状**: 日志显示权限错误或文件保存失败

**解决方案**:
```bash
# 检查 screenshots 目录权限
ls -la screenshots/

# 修复权限（Linux/Mac）
chmod 777 screenshots/

# 重启服务
docker-compose restart snapup-mcp
```

### 问题 4: Claude Desktop 无法连接

**症状**: Claude Desktop 报告 MCP 服务器连接失败

**检查清单**:
1. ✓ MCP 容器是否运行: `docker-compose ps snapup-mcp`
2. ✓ 配置文件路径是否正确
3. ✓ 使用 `-T` 参数（禁用伪终端）
4. ✓ 工作目录 (cwd) 是否指向正确位置

## 性能优化

### 减少启动延迟

默认配置在启动时等待 10 秒，你可以根据实际情况调整：

```yaml
# 修改 docker-compose.yml
command: sh -c "sleep 5 && /app/snapup -mode=mcp -output=/app/screenshots"
```

### 使用预构建镜像

```bash
# 先构建镜像
docker-compose build

# 后续启动更快
docker-compose --profile mcp up -d
```

## 更多信息

- [MCP 使用指南](./MCP_USAGE.md)
- [MCP 快速开始](./MCP_QUICKSTART.md)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [MCP 协议规范](https://spec.modelcontextprotocol.io/)

## 总结

通过 Docker 运行 MCP Server 的优势：

✅ 环境隔离，无需本地安装依赖  
✅ 与 Chrome 容器统一编排  
✅ 配置统一，易于部署  
✅ 支持多种运行模式（HTTP + MCP）  
✅ 使用 Docker Compose Profile 灵活控制启动  

推荐的工作流程：

1. 开发测试: 使用本地二进制文件，方便调试
2. 生产部署: 使用 Docker 容器，环境一致性更好
3. 混合模式: HTTP 服务用 Docker，MCP 用本地（根据需求）

