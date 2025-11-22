# SnapUp 中国部署指南

本指南专为中国用户优化，使用国内镜像源加速构建和部署。

## 🚀 快速开始（推荐）

使用中国版 Docker Compose，自动使用国内镜像源：

```bash
# 克隆项目
git clone <repository-url>
cd snapup

# 使用中国版配置构建并运行
make docker-run-cn

# 或直接使用 docker-compose
docker-compose -f docker-compose.cn.yml up -d
```

第一次构建可能需要 5-10 分钟，后续启动只需几秒钟。

访问 http://localhost:8080 开始使用！

## 📦 中国版优化内容

### 1. Go 依赖加速
使用七牛云 Go 代理镜像：
```bash
GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
```

### 2. Alpine Linux 镜像源
使用阿里云镜像：
```bash
mirrors.aliyun.com/alpine/
```

### 3. Debian 镜像源
使用中科大镜像：
```bash
mirrors.ustc.edu.cn/debian/
```

### 4. Chrome/Chromium
使用 Chromium 替代 Google Chrome，避免下载困难。

## 🔧 详细部署步骤

### 方法一：Docker Compose（推荐）

#### 1. 安装 Docker

**Ubuntu/Debian:**
```bash
# 使用阿里云 Docker 镜像
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

**CentOS/RHEL:**
```bash
# 使用阿里云 Docker 镜像
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo systemctl start docker
```

#### 2. 配置 Docker 镜像加速

创建或编辑 `/etc/docker/daemon.json`：

```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.ccs.tencentyun.com"
  ]
}
```

重启 Docker：
```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

#### 3. 构建并运行

```bash
cd snapup

# 使用中国版配置
make docker-run-cn

# 查看日志
docker-compose -f docker-compose.cn.yml logs -f

# 检查状态
docker-compose -f docker-compose.cn.yml ps
```

#### 4. 访问服务

浏览器打开：http://localhost:8080

### 方法二：本地构建

#### 1. 安装 Go

**使用官方安装包:**
```bash
# 下载 Go 1.21（使用国内镜像）
wget https://golang.google.cn/dl/go1.21.0.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc
```

**或使用包管理器:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang-1.21

# CentOS/RHEL
sudo yum install golang
```

#### 2. 配置 Go 代理

```bash
# 设置 Go 代理（七牛云）
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on

# 或使用阿里云
# go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

#### 3. 安装 Chromium

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install chromium-browser
```

**CentOS/RHEL:**
```bash
sudo yum install chromium
```

#### 4. 构建并运行

```bash
cd snapup

# 下载依赖
go mod download

# 构建
make build

# 运行（HTTP 模式）
./snapup -mode=http -port=8080

# 或运行（MCP 模式）
./snapup -mode=mcp -output=./screenshots
```

## 🐛 故障排除

### 问题 1: Docker 构建超时

**原因：** 网络连接问题

**解决方法：**
1. 确认已配置 Docker 镜像加速
2. 使用中国版 Dockerfile：
```bash
docker build -f Dockerfile.cn -t snapup:latest .
```

### 问题 2: Go 依赖下载失败

**解决方法：**
```bash
# 清理缓存
go clean -modcache

# 重新配置代理
go env -w GOPROXY=https://goproxy.cn,direct

# 重新下载
go mod download
```

### 问题 3: Chrome 下载失败

**解决方法：**
使用 Chromium 替代：
```bash
# Ubuntu/Debian
sudo apt-get install chromium-browser

# 或在 Dockerfile 中直接安装 Chromium
```

### 问题 4: Alpine 镜像下载慢

**解决方法：**
在 Dockerfile.cn 中已配置阿里云镜像：
```dockerfile
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
```

### 问题 5: 端口被占用

**解决方法：**
```bash
# 检查端口占用
sudo lsof -i :8080

# 修改端口
./snapup -mode=http -port=8088

# 或修改 docker-compose.cn.yml
```

## 📊 性能优化

### 1. 构建缓存

首次构建后，Docker 会缓存各层，后续构建会很快：
```bash
# 清理缓存（如需要）
docker builder prune

# 重新构建
make docker-cn
```

### 2. 资源限制

在 `docker-compose.cn.yml` 中设置资源限制：
```yaml
services:
  snapup:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### 3. 日志管理

限制日志大小：
```yaml
services:
  snapup:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 🔐 安全建议

### 1. 防火墙配置

```bash
# 仅允许本地访问
sudo ufw allow from 127.0.0.1 to any port 8080

# 或允许特定 IP
sudo ufw allow from 192.168.1.0/24 to any port 8080
```

### 2. 反向代理（Nginx）

```nginx
server {
    listen 80;
    server_name snapup.example.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. HTTPS 配置

使用 Let's Encrypt：
```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d snapup.example.com
```

## 📈 监控和维护

### 查看日志

```bash
# Docker 日志
docker-compose -f docker-compose.cn.yml logs -f snapup

# 仅查看最近 100 行
docker-compose -f docker-compose.cn.yml logs --tail=100 snapup
```

### 检查健康状态

```bash
# 检查服务状态
docker-compose -f docker-compose.cn.yml ps

# 检查健康检查
curl http://localhost:8080/api/health
```

### 重启服务

```bash
# 重启单个服务
docker-compose -f docker-compose.cn.yml restart snapup

# 重启所有服务
docker-compose -f docker-compose.cn.yml restart
```

### 更新服务

```bash
# 拉取最新代码
git pull

# 重新构建并启动
make docker-stop-cn
make docker-run-cn
```

## 🌐 网络配置

### 使用自定义端口

**方法 1: 修改 docker-compose.cn.yml**
```yaml
services:
  snapup:
    ports:
      - "8088:8080"  # 主机端口:容器端口
```

**方法 2: 环境变量**
```bash
PORT=8088 docker-compose -f docker-compose.cn.yml up -d
```

### 绑定特定 IP

```yaml
services:
  snapup:
    ports:
      - "127.0.0.1:8080:8080"  # 仅本地访问
```

## 🔄 数据备份

### 备份截图

```bash
# 手动备份
tar -czf screenshots-backup-$(date +%Y%m%d).tar.gz ./screenshots/

# 定时备份（添加到 crontab）
0 2 * * * cd /path/to/snapup && tar -czf /backup/screenshots-$(date +\%Y\%m\%d).tar.gz ./screenshots/
```

## 📚 相关资源

- **Go 中国镜像**: https://goproxy.cn
- **阿里云 Docker 镜像**: https://mirrors.aliyun.com/docker-ce/
- **中科大镜像站**: https://mirrors.ustc.edu.cn
- **腾讯云镜像**: https://mirrors.tencent.com

## 💡 最佳实践

1. **首次部署**: 使用 `docker-compose.cn.yml` 避免网络问题
2. **开发环境**: 本地构建，使用 Go 代理加速
3. **生产环境**: Docker 部署，配置反向代理和 HTTPS
4. **定期更新**: 每月检查并更新依赖和基础镜像
5. **监控日志**: 使用日志聚合工具（如 ELK）
6. **资源监控**: 使用 Prometheus + Grafana

## 🆘 获取帮助

如遇到问题：

1. 查看日志：`docker-compose -f docker-compose.cn.yml logs`
2. 检查网络：`ping mirrors.aliyun.com`
3. 验证配置：`docker-compose -f docker-compose.cn.yml config`
4. 提交 Issue: [GitHub Issues](https://github.com/gotoailab/snapup/issues)

## ✅ 验证安装

运行以下命令验证安装：

```bash
# 检查服务是否运行
curl http://localhost:8080/api/health

# 应返回
# {"status":"ok","service":"snapup"}

# 测试截图功能
curl -X POST http://localhost:8080/api/screenshot \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.baidu.com","device":"desktop","style":"none"}'
```

如果返回成功，说明服务正常运行！🎉

---

**祝你使用愉快！如有问题，欢迎反馈。**

