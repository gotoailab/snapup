# 安装指南

本文档提供 SnapUp 的详细安装说明。

## 目录

1. [Docker 安装（推荐）](#docker-安装推荐)
2. [本地安装](#本地安装)
3. [Chrome 安装](#chrome-安装)
4. [验证安装](#验证安装)

## Docker 安装（推荐）

使用 Docker 是最简单和可靠的安装方式。

### 前置要求

- Docker 20.10+
- Docker Compose 1.29+

### 安装步骤

1. **安装 Docker**

Ubuntu/Debian:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

2. **克隆并启动项目**

```bash
git clone <repository-url>
cd snapup
docker-compose up -d
```

3. **验证安装**

```bash
curl http://localhost:8080/api/health
# 期望输出: {"status":"ok","service":"snapup"}
```

## 本地安装

### 前置要求

- Go 1.21+
- Google Chrome

### 安装步骤

1. **安装 Go**

```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

2. **克隆和构建项目**

```bash
git clone <repository-url>
cd snapup
go mod download
go build -o snapup ./cmd/snapup
./snapup -port 8080
```

## Chrome 安装

### Ubuntu/Debian

```bash
# 添加 Google Chrome 仓库
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
sudo sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'

# 安装 Chrome
sudo apt-get update
sudo apt-get install google-chrome-stable

# 安装依赖
sudo apt-get install -y \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libgbm1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libx11-6 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2
```

### macOS

```bash
brew install --cask google-chrome
```

### 验证 Chrome 安装

```bash
google-chrome --version
```

## 验证安装

1. **检查服务状态**

```bash
curl http://localhost:8080/api/health
```

2. **测试截图功能**

```bash
curl -X POST http://localhost:8080/api/screenshot \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "device": "desktop",
    "style": "none"
  }'
```

3. **访问 Web 界面**

打开浏览器访问: http://localhost:8080

## 常见问题

### Chrome 无法启动

确保 Chrome 已正确安装并在系统 PATH 中：

```bash
which google-chrome
```

### 端口被占用

更改端口：

```bash
./snapup -port 8090
```

### 权限错误

确保输出目录有正确的权限：

```bash
mkdir -p screenshots
chmod 755 screenshots
```

## 生产环境部署

### 使用 systemd

```ini
# /etc/systemd/system/snapup.service
[Unit]
Description=SnapUp Screenshot Service
After=network.target

[Service]
Type=simple
User=snapup
WorkingDirectory=/opt/snapup
ExecStart=/opt/snapup/snapup -port 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl enable snapup
sudo systemctl start snapup
```

---

安装完成后，请参考 [README.md](README.md) 和 [QUICKSTART.md](QUICKSTART.md) 了解使用说明。
