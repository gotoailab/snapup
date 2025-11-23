# Chrome WebSocket 内网穿透 - 快速开始

这是一份 5 分钟快速配置指南。详细文档请参考：[完整配置指南](./docs/TUNNEL_SETUP.md)

## 🚀 方案选择

根据您的需求选择一种方案：

| 方案 | 适合场景 | 时间 |
|------|----------|------|
| **frp** | 有公网服务器 | 5 分钟 |
| **Cloudflare Tunnel** | 没有公网 IP，需要域名 | 10 分钟 |
| **Nginx + frp** | 需要认证保护 | 10 分钟 |

---

## 方案 1️⃣：使用 frp（最简单）

### 前置条件
- 一台有公网 IP 的服务器
- 在服务器上运行 frps（frp 服务端）

### 步骤 1：配置 frp 客户端

```bash
# 运行配置向导
./setup-tunnel.sh

# 或手动编辑 frpc.toml
nano frpc.toml
```

修改以下配置：
```toml
serverAddr = "your-server-ip.com"    # 改为您的服务器地址
serverPort = 7000                     # 服务器端口
auth.token = "your-secret-token"     # 改为您的密钥
```

### 步骤 2：启动服务

```bash
# 使用 Make 命令
make tunnel-frp

# 或直接使用 docker-compose
docker-compose -f docker-compose.tunnel.yml up -d
```

### 步骤 3：测试连接

```bash
# 使用测试脚本
./test-tunnel.sh http://your-server-ip:9222

# 或使用 Make 命令
make tunnel-test URL=http://your-server-ip:9222

# 或手动测试
curl http://your-server-ip:9222/json/version
```

### 步骤 4：在代码中使用

```go
// Go
allocCtx, _ := chromedp.NewRemoteAllocator(ctx, "ws://your-server-ip:9222")
```

```python
# Python
chrome_options.add_experimental_option("debuggerAddress", "your-server-ip:9222")
```

**完成！** ✅

---

## 方案 2️⃣：使用 Cloudflare Tunnel（免费）

### 前置条件
- Cloudflare 账号
- 一个托管在 Cloudflare 的域名

### 步骤 1：安装 cloudflared

```bash
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb
```

### 步骤 2：登录并创建 Tunnel

```bash
# 登录
cloudflared tunnel login

# 创建 tunnel
cloudflared tunnel create snapup-chrome

# 获取 token
cloudflared tunnel token snapup-chrome
```

### 步骤 3：配置 Docker Compose

编辑 `docker-compose.cloudflare.yml`，替换 `YOUR_TUNNEL_TOKEN`：

```yaml
command: tunnel --no-autoupdate run --token YOUR_TUNNEL_TOKEN
```

### 步骤 4：配置 DNS

```bash
cloudflared tunnel route dns snapup-chrome chrome.your-domain.com
```

### 步骤 5：启动服务

```bash
# 使用 Make 命令
make tunnel-cloudflare

# 或直接使用 docker-compose
docker-compose -f docker-compose.cloudflare.yml up -d
```

### 步骤 6：测试连接

```bash
# 测试 HTTPS 连接（Cloudflare 自动提供）
curl https://chrome.your-domain.com/json/version

# 或使用测试脚本
./test-tunnel.sh https://chrome.your-domain.com
```

### 步骤 7：在代码中使用

```go
// Go - 使用 WSS (安全 WebSocket)
allocCtx, _ := chromedp.NewRemoteAllocator(ctx, "wss://chrome.your-domain.com")
```

**完成！** ✅

---

## 方案 3️⃣：使用 Nginx（需要认证）

### 步骤 1：生成密码文件

```bash
# 安装 htpasswd 工具
sudo apt-get install apache2-utils

# 生成密码文件
htpasswd -c .htpasswd admin
```

### 步骤 2：更新 Nginx 配置

编辑 `nginx-websocket.conf`，添加认证：

```nginx
server {
    listen 9223;
    
    # 添加这两行
    auth_basic "Chrome DevTools";
    auth_basic_user_file /etc/nginx/.htpasswd;
    
    location / {
        # ... 其他配置
    }
}
```

### 步骤 3：更新 Docker Compose

编辑 `docker-compose.nginx-tunnel.yml`，在 nginx-ws 服务的 volumes 中添加：

```yaml
volumes:
  - ./nginx-websocket.conf:/etc/nginx/conf.d/default.conf:ro
  - ./.htpasswd:/etc/nginx/.htpasswd:ro  # 添加这行
```

### 步骤 4：配置 frpc.toml

```toml
serverAddr = "your-server-ip.com"
serverPort = 7000
auth.token = "your-secret-token"

[[proxies]]
name = "chrome-ws-nginx"
type = "tcp"
localIP = "nginx-ws"
localPort = 9223        # 注意是 Nginx 的端口
remotePort = 9223
```

### 步骤 5：启动服务

```bash
make tunnel-nginx
```

### 步骤 6：测试连接

```bash
# 带认证的测试
curl -u admin:password http://your-server-ip:9223/json/version
```

**完成！** ✅

---

## 🔍 故障排查

### 问题：无法连接

```bash
# 1. 检查容器状态
docker-compose ps

# 2. 查看日志
make tunnel-logs

# 3. 测试端口
nc -zv your-server-ip 9222

# 4. 检查防火墙
sudo ufw status
```

### 问题：frp 连接失败

```bash
# 查看 frpc 日志
docker logs snapup-frpc

# 检查 frps 服务端是否运行
# 在服务器上执行
ps aux | grep frps
```

### 问题：Cloudflare Tunnel 连接失败

```bash
# 查看 cloudflared 日志
docker logs snapup-cloudflared

# 检查 DNS 是否生效
nslookup chrome.your-domain.com
```

---

## 📚 常用命令

```bash
# 配置向导
make tunnel-setup

# 启动服务
make tunnel-frp              # frp 方案
make tunnel-cloudflare       # Cloudflare 方案
make tunnel-nginx            # Nginx 方案

# 查看日志
make tunnel-logs

# 测试连接
make tunnel-test URL=http://your-server-ip:9222

# 停止服务
make tunnel-stop

# 查看帮助
make help
```

---

## ⚠️ 安全提示

内网穿透会将 Chrome 暴露到公网，请务必：

1. **启用认证**：使用 HTTP Basic Auth 或更强的认证
2. **IP 白名单**：限制只允许特定 IP 访问
3. **使用 HTTPS**：通过 Cloudflare 或 SSL 证书加密传输
4. **监控日志**：定期检查访问日志
5. **定期更新**：及时更新 Chrome 和穿透工具

详细安全配置：[完整安全指南](./docs/TUNNEL_SETUP.md#安全建议)

---

## 🎯 下一步

1. ✅ 完成内网穿透配置
2. 📖 阅读[完整配置文档](./docs/TUNNEL_SETUP.md)
3. 🔐 配置[安全措施](./docs/TUNNEL_SETUP.md#安全建议)
4. 📊 监控性能和日志
5. 🚀 在您的项目中集成使用

---

## 💡 使用示例

### Go (chromedp)

```go
package main

import (
    "context"
    "github.com/chromedp/chromedp"
)

func main() {
    // 连接到远程 Chrome
    allocCtx, cancel := chromedp.NewRemoteAllocator(
        context.Background(),
        "ws://your-server-ip:9222",
    )
    defer cancel()

    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()

    var html string
    chromedp.Run(ctx,
        chromedp.Navigate("https://example.com"),
        chromedp.OuterHTML("html", &html),
    )
}
```

### Python (Selenium)

```python
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

chrome_options = Options()
chrome_options.add_experimental_option(
    "debuggerAddress", 
    "your-server-ip:9222"
)

driver = webdriver.Chrome(options=chrome_options)
driver.get("https://example.com")
```

### Node.js (Puppeteer)

```javascript
const puppeteer = require('puppeteer');

(async () => {
    const browser = await puppeteer.connect({
        browserWSEndpoint: 'ws://your-server-ip:9222'
    });
    
    const page = await browser.newPage();
    await page.goto('https://example.com');
})();
```

---

## 📞 获取帮助

- 📖 [完整文档](./docs/TUNNEL_SETUP.md)
- 🐛 [报告问题](https://github.com/your-repo/issues)
- 💬 [讨论区](https://github.com/your-repo/discussions)

---

**享受远程 Chrome 的便利！** 🎉

