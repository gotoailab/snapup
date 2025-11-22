# SnapUp Docker 部署指南

## 🚀 快速部署

### 中国用户（推荐）

使用国内镜像源加速：

```bash
cd snapup

# 构建并启动（使用国内镜像）
make docker-run-cn

# 或直接使用 docker-compose
docker-compose -f docker-compose.cn.yml up -d
```

### 国际用户

```bash
cd snapup

# 构建并启动
make docker-run

# 或直接使用 docker-compose
docker-compose up -d
```

## ✅ 验证部署

### 1. 检查容器状态

```bash
# 中国版
docker-compose -f docker-compose.cn.yml ps

# 国际版
docker-compose ps
```

预期输出：
```
NAME            STATUS                    PORTS
snapup-app      Up (healthy)             0.0.0.0:8080->8080/tcp
snapup-chrome   Up                       0.0.0.0:9222->9222/tcp
```

### 2. 测试健康接口

```bash
curl http://localhost:8080/api/health
```

预期输出：
```json
{"status":"ok","service":"snapup"}
```

### 3. 测试截图功能

```bash
curl -X POST http://localhost:8080/api/screenshot \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.baidu.com","device":"desktop","style":"none"}'
```

预期输出：
```json
{
  "success": true,
  "message": "截图成功",
  "image_url": "/screenshots/screenshot_desktop_xxx.png",
  "filename": "screenshot_desktop_xxx.png"
}
```

### 4. 访问 Web 界面

浏览器打开：http://localhost:8080

## 📝 常用命令

### 查看日志

```bash
# 所有服务日志
docker-compose -f docker-compose.cn.yml logs -f

# 仅查看 snapup 服务
docker-compose -f docker-compose.cn.yml logs -f snapup

# 仅查看 Chrome 服务
docker-compose -f docker-compose.cn.yml logs -f chrome

# 查看最近 50 行
docker-compose -f docker-compose.cn.yml logs --tail=50
```

### 重启服务

```bash
# 重启所有服务
docker-compose -f docker-compose.cn.yml restart

# 仅重启 snapup
docker-compose -f docker-compose.cn.yml restart snapup

# 仅重启 chrome
docker-compose -f docker-compose.cn.yml restart chrome
```

### 停止服务

```bash
# 停止但不删除容器
docker-compose -f docker-compose.cn.yml stop

# 停止并删除容器
docker-compose -f docker-compose.cn.yml down

# 停止并删除容器和数据卷
docker-compose -f docker-compose.cn.yml down -v
```

### 更新服务

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose -f docker-compose.cn.yml up -d --build
```

## 🔧 配置说明

### 端口配置

默认端口：
- HTTP 服务：8080
- Chrome 调试：9222

修改端口（编辑 `docker-compose.cn.yml`）：
```yaml
services:
  snapup:
    ports:
      - "8088:8080"  # 改为 8088
```

### 数据持久化

截图保存在本地 `./screenshots` 目录：
```yaml
volumes:
  - ./screenshots:/app/screenshots
```

### 资源限制

如需限制资源使用（编辑 `docker-compose.cn.yml`）：
```yaml
services:
  snapup:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

## 🐛 故障排除

### 问题 1: 容器启动失败

**检查日志：**
```bash
docker-compose -f docker-compose.cn.yml logs
```

**常见原因：**
- 端口被占用 → 修改端口号
- 磁盘空间不足 → 清理 Docker 镜像
- 权限问题 → 检查 screenshots 目录权限

### 问题 2: Chrome 容器无法启动

**检查 Chrome 日志：**
```bash
docker-compose -f docker-compose.cn.yml logs chrome
```

**解决方法：**
```bash
# 删除并重新创建
docker-compose -f docker-compose.cn.yml down
docker-compose -f docker-compose.cn.yml up -d
```

### 问题 3: 截图失败

**检查网络连接：**
```bash
# 测试容器间网络
docker exec snapup-app ping chrome -c 3
```

**检查 Chrome 状态：**
```bash
curl http://localhost:9222/json/version
```

### 问题 4: 构建速度慢（中国用户）

**解决方法：**
1. 使用中国版配置：`docker-compose.cn.yml`
2. 配置 Docker 镜像加速（见 INSTALL_CN.md）
3. 使用 `Dockerfile.cn` 构建

## 📊 性能优化

### 1. 构建缓存

Docker 会自动缓存构建层，加速后续构建：
```bash
# 清理缓存（如需要）
docker builder prune
```

### 2. 日志管理

限制日志大小（编辑 `docker-compose.cn.yml`）：
```yaml
services:
  snapup:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### 3. 定时清理截图

添加定时任务：
```bash
# 每天凌晨 2 点删除 7 天前的截图
0 2 * * * find /path/to/snapup/screenshots -name "*.png" -mtime +7 -delete
```

## 🔐 生产环境建议

### 1. 使用反向代理

**Nginx 示例：**
```nginx
server {
    listen 80;
    server_name snapup.example.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # 增加超时时间（截图可能需要较长时间）
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

### 2. 启用 HTTPS

```bash
sudo certbot --nginx -d snapup.example.com
```

### 3. 限制访问

**仅允许本地访问：**
```yaml
services:
  snapup:
    ports:
      - "127.0.0.1:8080:8080"
```

**或使用防火墙：**
```bash
sudo ufw allow from 192.168.1.0/24 to any port 8080
```

### 4. 监控和告警

使用 Docker 自带的监控：
```bash
docker stats snapup-app snapup-chrome
```

## 📈 架构说明

```
┌─────────────────┐
│   用户浏览器    │
└────────┬────────┘
         │ HTTP :8080
         ▼
┌─────────────────┐
│   snapup-app    │ ◄─┐
└────────┬────────┘   │
         │ WebSocket  │ Docker Network
         │ :9222      │ (snapup-network)
         ▼            │
┌─────────────────┐   │
│  snapup-chrome  │ ──┘
│ (headless-shell)│
└─────────────────┘
```

- **snapup-app**: 主应用容器，处理 HTTP 请求
- **snapup-chrome**: Chrome headless 容器，执行截图
- **snapup-network**: Docker 内部网络，连接两个容器

## 🆘 获取帮助

1. 查看日志： `docker-compose -f docker-compose.cn.yml logs`
2. 检查配置： `docker-compose -f docker-compose.cn.yml config`
3. 提交 Issue: [GitHub Issues](https://github.com/gotoailab/snapup/issues)

---

**部署完成！祝你使用愉快！** 🎉

