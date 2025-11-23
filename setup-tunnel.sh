#!/bin/bash
# SnapUp Chrome 内网穿透快速配置脚本

set -e

echo "======================================"
echo "  SnapUp Chrome 内网穿透配置向导"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 选择穿透方案
echo -e "${BLUE}请选择内网穿透方案:${NC}"
echo "1) frp (需要有公网服务器)"
echo "2) Cloudflare Tunnel (免费，推荐)"
echo "3) Nginx + frp (需要认证或高级配置)"
echo ""
read -p "请输入选项 (1-3): " choice

case $choice in
    1)
        echo -e "\n${GREEN}[frp 配置]${NC}"
        echo "需要填写以下信息:"
        
        # 获取服务器地址
        read -p "frp 服务器地址 (例: your-server.com): " server_addr
        read -p "frp 服务器端口 (默认: 7000): " server_port
        server_port=${server_port:-7000}
        
        # 获取认证 token
        read -p "认证 Token: " auth_token
        
        # 获取远程端口
        read -p "远程暴露端口 (默认: 9222): " remote_port
        remote_port=${remote_port:-9222}
        
        # 生成配置文件
        cat > frpc.toml <<EOF
# frp 客户端配置
serverAddr = "$server_addr"
serverPort = $server_port

# 认证信息
auth.token = "$auth_token"

# 穿透 Chrome WebSocket
[[proxies]]
name = "chrome-ws"
type = "tcp"
localIP = "chrome"
localPort = 9222
remotePort = $remote_port
EOF
        
        echo -e "\n${GREEN}✓ 配置文件已生成: frpc.toml${NC}"
        echo -e "\n启动命令:"
        echo -e "${YELLOW}docker-compose -f docker-compose.tunnel.yml up -d${NC}"
        echo -e "\n测试连接:"
        echo -e "${YELLOW}curl http://$server_addr:$remote_port/json/version${NC}"
        ;;
        
    2)
        echo -e "\n${GREEN}[Cloudflare Tunnel 配置]${NC}"
        echo ""
        echo "步骤 1: 安装 cloudflared CLI 工具"
        echo "-------------------------------------"
        echo "请先确保已安装 cloudflared，如果没有安装，执行:"
        echo -e "${YELLOW}wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb${NC}"
        echo -e "${YELLOW}sudo dpkg -i cloudflared-linux-amd64.deb${NC}"
        echo ""
        read -p "已安装 cloudflared? (y/n): " has_cloudflared
        
        if [[ "$has_cloudflared" != "y" ]]; then
            echo -e "${RED}请先安装 cloudflared 后再运行此脚本${NC}"
            exit 1
        fi
        
        echo -e "\n步骤 2: 登录 Cloudflare"
        echo "-------------------------------------"
        cloudflared tunnel login
        
        echo -e "\n步骤 3: 创建 Tunnel"
        echo "-------------------------------------"
        cloudflared tunnel create snapup-chrome
        
        echo -e "\n步骤 4: 获取 Tunnel Token"
        echo "-------------------------------------"
        echo "请执行以下命令获取 Token:"
        echo -e "${YELLOW}cloudflared tunnel token snapup-chrome${NC}"
        echo ""
        read -p "请粘贴 Tunnel Token: " tunnel_token
        
        # 更新 docker-compose 文件
        sed -i "s|YOUR_TUNNEL_TOKEN|$tunnel_token|g" docker-compose.cloudflare.yml
        
        echo -e "\n步骤 5: 配置域名"
        echo "-------------------------------------"
        read -p "您的域名 (例: chrome.example.com): " domain
        
        echo "配置 DNS..."
        cloudflared tunnel route dns snapup-chrome "$domain"
        
        echo -e "\n${GREEN}✓ Cloudflare Tunnel 配置完成!${NC}"
        echo -e "\n启动命令:"
        echo -e "${YELLOW}docker-compose -f docker-compose.cloudflare.yml up -d${NC}"
        echo -e "\n访问地址:"
        echo -e "${YELLOW}https://$domain/json/version${NC}"
        echo -e "\nWebSocket URL:"
        echo -e "${YELLOW}wss://$domain${NC}"
        ;;
        
    3)
        echo -e "\n${GREEN}[Nginx + frp 配置]${NC}"
        echo "这个方案会在 Chrome 前面添加 Nginx 反向代理"
        echo ""
        
        # frp 配置
        read -p "frp 服务器地址: " server_addr
        read -p "frp 服务器端口 (默认: 7000): " server_port
        server_port=${server_port:-7000}
        read -p "认证 Token: " auth_token
        read -p "远程暴露端口 (默认: 9223): " remote_port
        remote_port=${remote_port:-9223}
        
        # 生成 frpc 配置（针对 Nginx 端口）
        cat > frpc.toml <<EOF
# frp 客户端配置
serverAddr = "$server_addr"
serverPort = $server_port

auth.token = "$auth_token"

# 穿透 Nginx 代理端口
[[proxies]]
name = "chrome-ws-nginx"
type = "tcp"
localIP = "nginx-ws"
localPort = 9223
remotePort = $remote_port
EOF
        
        # 询问是否需要认证
        read -p "是否启用 HTTP 基础认证? (y/n): " enable_auth
        
        if [[ "$enable_auth" == "y" ]]; then
            read -p "用户名: " username
            read -sp "密码: " password
            echo ""
            
            # 生成 .htpasswd 文件
            if command -v htpasswd &> /dev/null; then
                echo "$password" | htpasswd -ci .htpasswd "$username"
                
                # 更新 Nginx 配置添加认证
                echo -e "\n${YELLOW}请手动编辑 nginx-websocket.conf 添加以下内容:${NC}"
                echo "auth_basic \"Chrome DevTools\";"
                echo "auth_basic_user_file /etc/nginx/.htpasswd;"
                echo ""
                echo "然后在 docker-compose.nginx-tunnel.yml 的 nginx-ws volumes 中添加:"
                echo "- ./.htpasswd:/etc/nginx/.htpasswd:ro"
            else
                echo -e "${YELLOW}未安装 htpasswd 工具，请手动创建 .htpasswd 文件${NC}"
            fi
        fi
        
        echo -e "\n${GREEN}✓ 配置文件已生成${NC}"
        echo -e "\n启动命令:"
        echo -e "${YELLOW}docker-compose -f docker-compose.nginx-tunnel.yml up -d${NC}"
        echo -e "\n测试连接:"
        echo -e "${YELLOW}curl http://$server_addr:$remote_port/json/version${NC}"
        ;;
        
    *)
        echo -e "${RED}无效的选项${NC}"
        exit 1
        ;;
esac

echo ""
echo "======================================"
echo -e "${GREEN}配置完成!${NC}"
echo "======================================"
echo ""
echo -e "${BLUE}下一步:${NC}"
echo "1. 检查生成的配置文件"
echo "2. 运行上面显示的启动命令"
echo "3. 测试连接是否正常"
echo ""
echo -e "${YELLOW}详细文档: docs/TUNNEL_SETUP.md${NC}"
echo ""
echo -e "${RED}安全提示:${NC}"
echo "- 建议配置防火墙规则"
echo "- 定期更换认证凭证"
echo "- 监控访问日志"
echo "- 使用 HTTPS/WSS 加密连接"
echo ""

