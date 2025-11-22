# Docker MCP 支持更新日志

**更新日期**: 2025-11-22

## 概述

本次更新为 SnapUp 添加了完整的 Docker MCP Server 支持，使用户可以通过 Docker 容器运行 MCP 服务，实现更好的环境隔离和部署便利性。

## 主要变更

### 1. Dockerfile 增强

**文件**: `Dockerfile`

**变更内容**:
- 添加构建参数 `GOPROXY` 支持，可在构建时指定 Go 代理
- 添加 `GO111MODULE` 环境变量
- 增强灵活性，支持不同网络环境下的构建

**新增代码**:
```dockerfile
# 构建参数
ARG GOPROXY=""

# 如果提供了 GOPROXY，则设置 Go 代理
ENV GOPROXY=${GOPROXY}
ENV GO111MODULE=on
```

**用途**:
- 国内用户可通过 `--build-arg GOPROXY=https://goproxy.cn,direct` 加速构建
- 企业用户可指定内部 Go 代理
- 保持向后兼容，默认不设置代理

### 2. Docker Compose 配置更新

**文件**: `docker-compose.yml`

**变更内容**:
- 添加 `snapup-mcp` 服务配置
- 使用 Docker Compose Profile 特性，默认不启动 MCP 服务
- 配置 stdin/tty 支持 MCP 协议通信

**新增服务**:
```yaml
snapup-mcp:
  build:
    context: .
    dockerfile: Dockerfile
  container_name: snapup-mcp
  restart: "no"  # 不自动重启，需要手动启动
  environment:
    - RUN_MODE=mcp
    - OUTPUT_DIR=/app/screenshots
    - CHROME_WS_URL=ws://chrome:9222
  volumes:
    - ./screenshots:/app/screenshots
  depends_on:
    - chrome
  command: sh -c "sleep 10 && /app/snapup -mode=mcp -output=/app/screenshots"
  networks:
    - snapup-network
  stdin_open: true  # 保持 stdin 打开，MCP 协议需要
  tty: true         # 分配伪终端
  profiles:
    - mcp           # 使用 profile，默认不启动
```

**特点**:
- 使用 Profile 功能，需显式指定 `--profile mcp` 才启动
- 与 Chrome 容器共享网络
- 挂载 screenshots 目录用于保存截图
- 配置适合 stdio 通信的容器环境

### 3. MCP Docker 管理脚本

**文件**: `scripts/mcp-docker.sh`（新增）

**功能**:
- 一站式管理 Docker MCP 服务
- 支持启动、停止、重启、查看日志等操作
- 自动检测 Docker 环境
- 支持标准版和中国版配置文件切换
- 彩色输出和友好的用户提示

**主要命令**:
```bash
./scripts/mcp-docker.sh start        # 启动服务
./scripts/mcp-docker.sh stop         # 停止服务
./scripts/mcp-docker.sh restart      # 重启服务
./scripts/mcp-docker.sh logs -f      # 查看日志
./scripts/mcp-docker.sh status       # 查看状态
./scripts/mcp-docker.sh build        # 重新构建
./scripts/mcp-docker.sh start --cn   # 使用国内镜像源
```

**特色功能**:
- ✅ 自动检查 Docker 和 docker-compose 是否安装
- ✅ 自动检查 Chrome 服务是否运行
- ✅ 彩色输出，区分信息、成功、警告、错误
- ✅ 支持国内镜像源切换
- ✅ 友好的帮助信息

### 4. Makefile 更新

**文件**: `Makefile`

**新增命令**:
- `make docker-mcp-start`: 启动 Docker MCP 服务
- `make docker-mcp-stop`: 停止 Docker MCP 服务
- `make docker-mcp-restart`: 重启 Docker MCP 服务
- `make docker-mcp-logs`: 查看 Docker MCP 服务日志
- `make docker-mcp-status`: 查看 Docker MCP 服务状态
- `make docker-mcp-build`: 重新构建 Docker MCP 服务
- `make docker-mcp-start-cn`: 启动 Docker MCP 服务（中国版）

**更新内容**:
- 增强 help 命令，分类显示所有可用命令
- 添加 Docker MCP 命令分类
- 更新 .PHONY 声明

**使用示例**:
```bash
make docker-mcp-start    # 最简单的启动方式
make docker-mcp-logs     # 查看实时日志
make docker-mcp-status   # 检查服务状态
```

### 5. 文档更新

#### 5.1 Docker MCP 详细文档

**文件**: `docs/DOCKER_MCP.md`（新增）

**内容**:
- Docker MCP Server 完整使用指南
- 快速开始教程
- 配置说明
- 使用场景示例
- Claude Desktop 集成配置
- 常见命令参考
- 故障排除指南
- 性能优化建议

**覆盖场景**:
1. 同时运行 HTTP 和 MCP 模式
2. 仅运行 MCP 模式
3. 临时测试 MCP 服务
4. 在 Claude Desktop 中配置
5. 开发调试
6. 生产部署

#### 5.2 Docker MCP 快速示例

**文件**: `examples/docker-mcp-quickstart.md`（新增）

**内容**:
- 7 个实用场景示例
- 三种操作方式对比（Make / 脚本 / Docker Compose）
- 故障排除步骤
- 常用命令速查表
- 最佳实践

**示例场景**:
1. 快速启动 MCP 服务
2. 中国用户使用国内镜像源
3. 配置 Claude Desktop 使用 Docker MCP
4. 同时运行 HTTP 和 MCP 服务
5. 开发调试
6. 生产部署
7. 故障排除

#### 5.3 README 更新

**文件**: `README.md`

**更新内容**:
- 在 MCP 模式章节添加 Docker 运行说明
- 添加 Docker 模式下的 Claude Desktop 配置示例
- 添加新文档的链接引用
- 更新使用说明，区分本地运行和 Docker 运行

## 使用方式对比

### 方式 1: Make 命令（推荐）

**优点**: 最简单，记忆方便
**缺点**: 需要 Make 工具

```bash
make docker-mcp-start
make docker-mcp-logs
make docker-mcp-stop
```

### 方式 2: 脚本命令

**优点**: 功能丰富，有交互式提示
**缺点**: 命令稍长

```bash
./scripts/mcp-docker.sh start
./scripts/mcp-docker.sh logs -f
./scripts/mcp-docker.sh stop
```

### 方式 3: Docker Compose 原生命令

**优点**: 标准 Docker 工具，灵活性高
**缺点**: 命令较长，需要记住 profile 参数

```bash
docker-compose --profile mcp up -d
docker-compose logs -f snapup-mcp
docker-compose stop snapup-mcp
```

## 技术亮点

### 1. Docker Compose Profile

使用 Profile 特性实现可选服务：
- HTTP 服务默认启动（常规 Web 应用）
- MCP 服务按需启动（AI 工具集成）
- 避免资源浪费，提高灵活性

### 2. 构建参数支持

通过 ARG 指令支持构建时参数：
- 适应不同网络环境
- 企业内部部署友好
- 保持向后兼容

### 3. stdin/tty 配置

正确配置容器通信方式：
- `stdin_open: true` - 保持标准输入打开
- `tty: true` - 分配伪终端
- 支持 MCP 协议的 stdio 通信模式

### 4. 服务依赖管理

合理配置服务依赖：
- MCP 依赖 Chrome 服务
- 启动延迟确保 Chrome 就绪
- 健康检查保证服务质量

### 5. 管理脚本

提供友好的管理脚本：
- 自动环境检查
- 彩色输出提示
- 国内镜像源支持
- 完整的帮助信息

## 兼容性

### 向后兼容

✅ 现有 HTTP 服务不受影响  
✅ 现有 docker-compose.yml 配置完全兼容  
✅ 现有 Makefile 命令全部保留  
✅ 本地运行方式不受影响  

### 新增功能

✨ Docker 运行 MCP 服务  
✨ 灵活的服务启动控制  
✨ 丰富的管理工具  
✨ 完整的文档支持  

## 测试验证

### 构建测试

```bash
# 标准构建
docker-compose build snapup-mcp

# 使用 GOPROXY 构建
docker-compose build --build-arg GOPROXY=https://goproxy.cn,direct snapup-mcp

# 中国版构建
docker-compose -f docker-compose.cn.yml build snapup-mcp
```

### 功能测试

```bash
# 启动测试
make docker-mcp-start
make docker-mcp-status

# 日志测试
make docker-mcp-logs

# 重启测试
make docker-mcp-restart

# 停止测试
make docker-mcp-stop
```

### 集成测试

```bash
# 1. 启动服务
make docker-mcp-start

# 2. 测试容器通信
docker exec -i snapup-mcp /app/snapup -mode=mcp -output=/app/screenshots

# 3. 测试 Chrome 连接
docker exec snapup-mcp wget -O- http://chrome:9222/json/version

# 4. 测试截图保存
docker exec snapup-mcp ls -la /app/screenshots/
```

## 部署建议

### 开发环境

- 使用本地运行方式，方便调试
- 使用 `make run-mcp` 快速启动

### 测试环境

- 使用 Docker 运行，环境一致
- 使用 `make docker-mcp-start` 启动服务

### 生产环境

- 使用 Docker Compose 部署
- 配置适当的重启策略
- 配置日志轮转
- 配置资源限制

## 后续计划

### 短期改进

- [ ] 添加健康检查端点
- [ ] 支持更多环境变量配置
- [ ] 添加 Prometheus 监控指标

### 长期规划

- [ ] Kubernetes 部署支持
- [ ] 多实例负载均衡
- [ ] 分布式截图任务队列

## 总结

本次更新为 SnapUp 添加了完整的 Docker MCP Server 支持，主要优势：

✅ **灵活部署**: 支持本地运行和 Docker 运行  
✅ **环境隔离**: Docker 容器提供一致的运行环境  
✅ **易于管理**: 提供多种管理工具（Make / 脚本 / Docker Compose）  
✅ **国际友好**: 支持国内镜像源加速  
✅ **文档完善**: 提供详细的使用文档和示例  
✅ **向后兼容**: 不影响现有功能和使用方式  

通过这些改进，用户现在可以：
1. 通过 Docker 快速部署 MCP Server
2. 在 Claude Desktop 中使用 Docker 容器化的 MCP 服务
3. 灵活选择运行方式（本地 vs Docker）
4. 使用友好的管理工具控制服务
5. 参考完整的文档和示例快速上手

## 相关文档

- [Docker MCP 详细文档](docs/DOCKER_MCP.md)
- [Docker MCP 快速示例](examples/docker-mcp-quickstart.md)
- [MCP 使用指南](docs/MCP_USAGE.md)
- [项目 README](README.md)

