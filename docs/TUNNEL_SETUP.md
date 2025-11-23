# Chrome WebSocket 内网穿透配置指南

本文档介绍如何将 Chrome 的 WebSocket 协议（端口 9222）通过内网穿透暴露到公网，方便远程调用。

## 为什么需要内网穿透？

在某些场景下，您可能需要从公网访问 Chrome DevTools Protocol：

- 在云端运行爬虫或自动化脚本，连接到本地 Chrome
- 多人共享同一个 Chrome 实例进行调试
- 在不同网络环境下远程控制 Chrome

## 方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|---------|
| **frp** | 稳定、快速、可自建服务端 | 需要有公网服务器 | 企业内部使用 |
| **Cloudflare Tunnel** | 免费、无需公网 IP、自动 HTTPS | 速度依赖 Cloudflare 网络 | 个人项目、国外用户 |
| **Nginx + 穿透** | 灵活、可添加认证 | 配置较复杂 | 需要精细控制 |

---

## 方案一：使用 frp（推荐）

### 前置条件

1. 一台有公网 IP 的服务器（作为 frp 服务端）
2. 在服务器上安装并运行 frps（frp 服务端）

### 步骤 1：配置 frp 服务端

在您的公网服务器上创建 `frps.toml`：

```toml
# frp 服务端配置
bindPort = 7000
auth.token = "your-secret-token-here"  # 请修改为强密码

# 允许的端口范围
allowPorts = [
  { start = 9000, end = 9999 }
]

# Web 控制台（可选）
webServer.port = 7500
webServer.user = "admin"
webServer.password = "admin123"
```

启动 frps：

```bash
# 下载 frp（在服务器上）
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64

# 启动服务端
./frps -c frps.toml

# 或使用 systemd 服务
sudo nano /etc/systemd/system/frps.service
```

systemd 服务文件内容：

```ini
[Unit]
Description=frp server
After=network.target

[Service]
Type=simple
ExecStart=/path/to/frps -c /path/to/frps.toml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable frps
sudo systemctl start frps
sudo systemctl status frps
```

### 步骤 2：配置 frp 客户端

编辑 `frpc.toml`（项目根目录已创建）：

```toml
serverAddr = "your-server-ip.com"    # 替换为您的服务器 IP 或域名
serverPort = 7000
auth.token = "your-secret-token-here"  # 与服务端保持一致

# TCP 模式穿透
[[proxies]]
name = "chrome-ws"
type = "tcp"
localIP = "chrome"
localPort = 9222
remotePort = 9222

# 或使用 HTTP 模式（如果有域名）
[[proxies]]
name = "chrome-ws-http"
type = "http"
localIP = "chrome"
localPort = 9222
customDomains = ["chrome.your-domain.com"]
```

### 步骤 3：启动服务

```bash
# 使用带 frp 的 docker-compose
docker-compose -f docker-compose.tunnel.yml up -d

# 查看日志
docker-compose -f docker-compose.tunnel.yml logs -f frpc
```

### 步骤 4：测试连接

```bash
# 测试 WebSocket 连接
curl http://your-server-ip:9222/json/version

# 或者使用域名（如果配置了 HTTP 模式）
curl http://chrome.your-domain.com/json/version
```

### 步骤 5：在代码中使用

修改您的应用配置，使用公网地址：

```bash
# 环境变量
export CHROME_WS_URL="ws://your-server-ip:9222"

# 或使用域名
export CHROME_WS_URL="ws://chrome.your-domain.com"
```

---

## 方案二：使用 Cloudflare Tunnel（免费）

### 优点
- 完全免费
- 不需要公网 IP
- 自动 HTTPS
- 稳定可靠

### 步骤 1：安装 cloudflared

```bash
# 在本地安装 cloudflared CLI（用于配置）
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb

# 登录 Cloudflare
cloudflared tunnel login
```

### 步骤 2：创建 Tunnel

```bash
# 创建 tunnel
cloudflared tunnel create snapup-chrome

# 记录返回的 Tunnel ID 和 credentials 文件路径
# 例如: Tunnel ID: 12345678-1234-1234-1234-123456789abc
```

### 步骤 3：配置 Tunnel

创建 `cloudflared-config.yml`：

```yaml
tunnel: 12345678-1234-1234-1234-123456789abc  # 替换为您的 Tunnel ID
credentials-file: /etc/cloudflared/credentials.json

ingress:
  # Chrome WebSocket
  - hostname: chrome.your-domain.com  # 替换为您的域名
    service: http://chrome:9222
    originRequest:
      noTLSVerify: true
  
  # Catch-all rule
  - service: http_status:404
```

### 步骤 4：配置 DNS

在 Cloudflare Dashboard 中添加 DNS 记录：

```bash
# 使用命令行添加
cloudflared tunnel route dns snapup-chrome chrome.your-domain.com
```

### 步骤 5：使用 Docker 运行

编辑 `docker-compose.cloudflare.yml`，替换 `YOUR_TUNNEL_TOKEN`：

```bash
# 获取 tunnel token
cloudflared tunnel token snapup-chrome

# 启动服务
docker-compose -f docker-compose.cloudflare.yml up -d

# 查看日志
docker-compose -f docker-compose.cloudflare.yml logs -f cloudflared
```

### 步骤 6：测试访问

```bash
# 测试连接
curl https://chrome.your-domain.com/json/version

# 使用 WebSocket
export CHROME_WS_URL="wss://chrome.your-domain.com"
```

---

## 方案三：使用 Nginx + 任意穿透工具

这个方案的优势是在 Chrome 前面加了一层 Nginx，可以：

1. 添加基础认证
2. 添加访问日志
3. 限流保护
4. SSL 终端

### 步骤 1：使用 Nginx 代理

已经创建好了 `nginx-websocket.conf` 配置文件。

### 步骤 2：添加认证（可选）

编辑 `nginx-websocket.conf`，添加基础认证：

```nginx
server {
    listen 9223;
    server_name _;

    # 基础认证
    auth_basic "Chrome DevTools";
    auth_basic_user_file /etc/nginx/.htpasswd;

    location / {
        # ... 其他配置保持不变
    }
}
```

生成密码文件：

```bash
# 安装 htpasswd 工具
sudo apt-get install apache2-utils

# 生成密码文件
htpasswd -c .htpasswd admin

# 将密码文件挂载到容器
# 在 docker-compose 中添加 volume
volumes:
  - ./nginx-websocket.conf:/etc/nginx/conf.d/default.conf:ro
  - ./.htpasswd:/etc/nginx/.htpasswd:ro
```

### 步骤 3：启动服务

```bash
# 启动 Nginx + frp
docker-compose -f docker-compose.nginx-tunnel.yml up -d

# 查看日志
docker-compose -f docker-compose.nginx-tunnel.yml logs -f
```

### 步骤 4：测试

```bash
# 不带认证（如果没配置）
curl http://your-server-ip:9223/json/version

# 带认证
curl -u admin:password http://your-server-ip:9223/json/version
```

---

## 安全建议

⚠️ **重要：将 Chrome DevTools 暴露到公网存在安全风险！**

### 1. 使用认证

在 Nginx 中添加基础认证或更复杂的认证机制。

### 2. IP 白名单

限制只允许特定 IP 访问：

```nginx
# 在 nginx-websocket.conf 中添加
location / {
    allow 1.2.3.4;      # 允许的 IP
    allow 5.6.7.0/24;   # 允许的 IP 段
    deny all;           # 拒绝其他所有 IP
    
    proxy_pass http://chrome_backend;
    # ... 其他配置
}
```

### 3. 使用 HTTPS/WSS

通过 Cloudflare Tunnel 或在 Nginx 中配置 SSL：

```nginx
server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # ... 其他配置
}
```

### 4. 限流保护

防止滥用：

```nginx
# 在 nginx-websocket.conf 中添加
limit_req_zone $binary_remote_addr zone=chrome_limit:10m rate=10r/s;

server {
    location / {
        limit_req zone=chrome_limit burst=20 nodelay;
        # ... 其他配置
    }
}
```

### 5. 防火墙规则

在服务器上配置防火墙：

```bash
# 只允许特定 IP 访问
sudo ufw allow from 1.2.3.4 to any port 9222
sudo ufw deny 9222
```

---

## 高级用法

### 1. 连接到公网 Chrome

从任意地方连接到您的 Chrome：

```go
// Go 代码示例
import "github.com/chromedp/chromedp"

func main() {
    // 连接到公网 Chrome
    allocCtx, cancel := chromedp.NewRemoteAllocator(
        context.Background(),
        "ws://your-server-ip:9222",
    )
    defer cancel()

    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()

    // 使用 Chrome
    var html string
    err := chromedp.Run(ctx,
        chromedp.Navigate("https://example.com"),
        chromedp.OuterHTML("html", &html),
    )
}
```

```python
# Python 代码示例
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

chrome_options = Options()
chrome_options.add_experimental_option("debuggerAddress", "your-server-ip:9222")

driver = webdriver.Chrome(options=chrome_options)
driver.get("https://example.com")
```

### 2. 多实例负载均衡

如果需要运行多个 Chrome 实例：

```yaml
# docker-compose 配置
services:
  chrome1:
    # ... Chrome 配置
    ports:
      - "9222:9222"
  
  chrome2:
    # ... Chrome 配置
    ports:
      - "9223:9222"
  
  nginx-lb:
    # Nginx 负载均衡配置
    # ...
```

---

## 故障排查

### 问题 1：WebSocket 连接失败

**症状**：无法建立 WebSocket 连接

**解决方案**：

1. 检查防火墙是否开放端口
2. 确认 Nginx 配置中有 `Upgrade` 和 `Connection` 头
3. 检查 frp 日志是否有错误

```bash
# 测试端口是否开放
nc -zv your-server-ip 9222

# 查看 frp 日志
docker-compose logs frpc
```

### 问题 2：Chrome 容器无法启动

**症状**：Chrome 容器反复重启

**解决方案**：

1. 增加 shm_size
2. 检查 Chrome 启动参数

```bash
# 查看 Chrome 日志
docker logs snapup-chrome

# 手动测试 Chrome
docker exec -it snapup-chrome /headless-shell/headless-shell --version
```

### 问题 3：连接超时

**症状**：长时间运行后连接断开

**解决方案**：

调整超时设置：

```nginx
# 在 nginx-websocket.conf 中
proxy_read_timeout 7d;
proxy_send_timeout 7d;
```

```toml
# 在 frpc.toml 中
[common]
heartbeat_interval = 30
heartbeat_timeout = 90
```

---

## 性能优化

### 1. 使用本地 Chrome 缓存

```bash
# 挂载 Chrome 用户数据目录
docker run -v chrome-data:/home/chrome/.config/google-chrome snapup-chrome
```

### 2. 启用 Chrome 压缩

```bash
# 在 Dockerfile.chrome 中添加
--enable-features=NetworkService,NetworkServiceInProcess
```

### 3. 调整 Nginx 缓冲

```nginx
proxy_buffering off;
proxy_request_buffering off;
```

---

## 参考资料

- [frp 官方文档](https://github.com/fatedier/frp)
- [Cloudflare Tunnel 文档](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/)
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)
- [Nginx WebSocket 代理](https://nginx.org/en/docs/http/websocket.html)

---

## 总结

根据您的使用场景选择合适的方案：

- **企业内部使用** → 选择 frp，自建服务端更安全可控
- **个人项目/快速部署** → 选择 Cloudflare Tunnel，免费且简单
- **需要精细控制** → 选择 Nginx + 穿透，可以添加各种中间件

无论选择哪种方案，都要注意安全性，不要将未经保护的 Chrome 暴露到公网！

