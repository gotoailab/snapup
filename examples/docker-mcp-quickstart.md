# Docker MCP 快速开始示例

本文档提供使用 Docker 运行 SnapUp MCP Server 的快速开始示例。

## 场景 1: 快速启动 MCP 服务

### 使用 Make 命令（最简单）

```bash
# 启动 MCP 服务
make docker-mcp-start

# 查看服务状态
make docker-mcp-status

# 查看实时日志
make docker-mcp-logs

# 停止服务
make docker-mcp-stop
```

### 使用脚本命令

```bash
# 启动服务
./scripts/mcp-docker.sh start

# 查看状态
./scripts/mcp-docker.sh status

# 查看日志
./scripts/mcp-docker.sh logs -f

# 停止服务
./scripts/mcp-docker.sh stop
```

### 使用 Docker Compose 命令

```bash
# 启动 MCP 服务
docker-compose --profile mcp up -d

# 查看日志
docker-compose logs -f snapup-mcp

# 停止服务
docker-compose stop snapup-mcp
```

## 场景 2: 中国用户使用国内镜像源

### 使用 Make 命令

```bash
# 使用国内镜像源启动
make docker-mcp-start-cn

# 查看状态
make docker-mcp-status
```

### 使用脚本命令

```bash
# 使用国内镜像源启动
./scripts/mcp-docker.sh start --cn

# 查看日志
./scripts/mcp-docker.sh logs -f --cn
```

### 使用 Docker Compose 命令

```bash
# 使用中国版配置文件启动
docker-compose -f docker-compose.cn.yml --profile mcp up -d

# 查看日志
docker-compose -f docker-compose.cn.yml logs -f snapup-mcp
```

## 场景 3: 配置 Claude Desktop 使用 Docker MCP

### 步骤 1: 启动 MCP 服务

```bash
# 确保服务正在运行
make docker-mcp-start
make docker-mcp-status
```

### 步骤 2: 配置 Claude Desktop

编辑配置文件：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

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
      "cwd": "/absolute/path/to/snapup"
    }
  }
}
```

**注意**: 
- 将 `/absolute/path/to/snapup` 替换为你的 SnapUp 项目的绝对路径
- `-T` 参数很重要，用于禁用伪终端分配

### 步骤 3: 重启 Claude Desktop

完全退出 Claude Desktop，然后重新启动。

### 步骤 4: 测试

在 Claude Desktop 中说：

```
请帮我截取 https://www.google.com 的桌面版截图
```

## 场景 4: 同时运行 HTTP 和 MCP 服务

```bash
# 启动 HTTP 服务（默认）
docker-compose up -d

# 同时启动 MCP 服务
docker-compose --profile mcp up -d snapup-mcp

# 查看所有服务状态
docker-compose ps

# 此时可以：
# - 通过浏览器访问 http://localhost:8080 使用 HTTP 接口
# - 通过 Claude Desktop 使用 MCP 接口
```

## 场景 5: 开发调试

### 查看详细日志

```bash
# 实时查看 MCP 服务日志
make docker-mcp-logs

# 查看所有服务日志
docker-compose logs -f

# 查看最近 50 行日志
docker-compose logs --tail=50 snapup-mcp
```

### 重新构建服务

```bash
# 代码更改后重新构建
make docker-mcp-build

# 重启服务
make docker-mcp-restart

# 或者一次性完成
docker-compose --profile mcp up -d --build snapup-mcp
```

### 进入容器调试

```bash
# 进入 MCP 容器
docker exec -it snapup-mcp /bin/sh

# 查看文件
ls -la /app/

# 查看环境变量
env | grep -E 'CHROME|OUTPUT'

# 测试网络连接
wget -O- http://chrome:9222/json/version

# 退出容器
exit
```

## 场景 6: 生产部署

### 准备工作

```bash
# 构建优化的生产镜像
docker-compose build --no-cache

# 启动服务
docker-compose up -d chrome
docker-compose --profile mcp up -d snapup-mcp

# 验证服务
make docker-mcp-status
```

### 配置自动重启（可选）

如果你希望 MCP 服务在失败时自动重启，编辑 `docker-compose.yml`:

```yaml
snapup-mcp:
  # ... 其他配置 ...
  restart: "unless-stopped"  # 改为 unless-stopped
```

然后重新启动：

```bash
docker-compose --profile mcp up -d snapup-mcp
```

### 监控和维护

```bash
# 定期检查服务状态
watch -n 10 'docker-compose ps'

# 查看资源使用
docker stats snapup-mcp

# 清理旧日志（如果日志太大）
docker-compose logs --tail=0 snapup-mcp > /dev/null
```

## 场景 7: 故障排除

### MCP 服务无法启动

```bash
# 1. 检查 Chrome 服务是否运行
docker-compose ps chrome

# 2. 如果 Chrome 未运行，先启动它
docker-compose up -d chrome
sleep 10

# 3. 再启动 MCP 服务
make docker-mcp-start

# 4. 查看详细日志
make docker-mcp-logs
```

### 无法连接到 Chrome

```bash
# 1. 测试 Chrome 是否可访问
docker exec snapup-mcp wget -O- http://chrome:9222/json/version

# 2. 如果失败，重启 Chrome
docker-compose restart chrome
sleep 10

# 3. 重启 MCP 服务
make docker-mcp-restart
```

### 截图保存失败

```bash
# 1. 检查 screenshots 目录权限
ls -la screenshots/

# 2. 修复权限
chmod 777 screenshots/

# 3. 检查容器内的挂载
docker exec snapup-mcp ls -la /app/screenshots/

# 4. 重启服务
make docker-mcp-restart
```

### Claude Desktop 无法连接

```bash
# 1. 确保 MCP 容器正在运行
docker-compose ps snapup-mcp

# 2. 测试容器交互
docker exec -i snapup-mcp /app/snapup -mode=mcp -output=/app/screenshots

# 3. 检查 Claude Desktop 配置
# 确保使用了 -T 参数和正确的 cwd
```

## 常用命令速查

| 操作 | Make 命令 | 脚本命令 | Docker Compose 命令 |
|------|----------|---------|-------------------|
| 启动服务 | `make docker-mcp-start` | `./scripts/mcp-docker.sh start` | `docker-compose --profile mcp up -d` |
| 停止服务 | `make docker-mcp-stop` | `./scripts/mcp-docker.sh stop` | `docker-compose stop snapup-mcp` |
| 重启服务 | `make docker-mcp-restart` | `./scripts/mcp-docker.sh restart` | `docker-compose restart snapup-mcp` |
| 查看日志 | `make docker-mcp-logs` | `./scripts/mcp-docker.sh logs -f` | `docker-compose logs -f snapup-mcp` |
| 查看状态 | `make docker-mcp-status` | `./scripts/mcp-docker.sh status` | `docker-compose ps` |
| 重新构建 | `make docker-mcp-build` | `./scripts/mcp-docker.sh build` | `docker-compose build snapup-mcp` |

## 更多资源

- [Docker MCP 详细文档](../docs/DOCKER_MCP.md)
- [MCP 使用指南](../docs/MCP_USAGE.md)
- [MCP 快速开始](../docs/MCP_QUICKSTART.md)
- [项目 README](../README.md)

