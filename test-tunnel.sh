#!/bin/bash
# SnapUp Chrome 内网穿透测试脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "======================================"
echo "  SnapUp Chrome 内网穿透测试"
echo "======================================"
echo ""

# 检查参数
if [ -z "$1" ]; then
    echo -e "${YELLOW}使用方法:${NC}"
    echo "  $0 <测试地址>"
    echo ""
    echo "示例:"
    echo "  $0 http://your-server-ip:9222"
    echo "  $0 http://chrome.your-domain.com"
    echo "  $0 https://chrome.your-domain.com"
    echo ""
    exit 1
fi

TEST_URL="$1"
echo -e "${BLUE}测试地址:${NC} $TEST_URL"
echo ""

# 测试函数
test_endpoint() {
    local endpoint=$1
    local description=$2
    local full_url="${TEST_URL}${endpoint}"
    
    echo -e "${BLUE}测试:${NC} $description"
    echo -e "URL: $full_url"
    
    # 使用 curl 测试
    response=$(curl -s -w "\nHTTP_CODE:%{http_code}\n" "$full_url" 2>&1)
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d':' -f2)
    content=$(echo "$response" | grep -v "HTTP_CODE:")
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓ 成功 (HTTP $http_code)${NC}"
        if [ ! -z "$content" ]; then
            echo -e "${GREEN}响应内容:${NC}"
            echo "$content" | head -5
        fi
        return 0
    else
        echo -e "${RED}✗ 失败 (HTTP $http_code)${NC}"
        if [ ! -z "$content" ]; then
            echo -e "${RED}错误信息:${NC}"
            echo "$content" | head -5
        fi
        return 1
    fi
    echo ""
}

# 执行测试
total_tests=0
passed_tests=0

# 测试 1: Chrome 版本信息
echo "======================================"
echo "测试 1: Chrome 版本信息"
echo "======================================"
if test_endpoint "/json/version" "获取 Chrome 版本"; then
    ((passed_tests++))
fi
((total_tests++))
echo ""

# 测试 2: 浏览器目标列表
echo "======================================"
echo "测试 2: 浏览器目标列表"
echo "======================================"
if test_endpoint "/json/list" "获取浏览器目标"; then
    ((passed_tests++))
fi
((total_tests++))
echo ""

# 测试 3: 协议信息
echo "======================================"
echo "测试 3: DevTools 协议"
echo "======================================"
if test_endpoint "/json/protocol" "获取协议信息"; then
    ((passed_tests++))
fi
((total_tests++))
echo ""

# 测试 4: 新建页面
echo "======================================"
echo "测试 4: 创建新页面"
echo "======================================"
if test_endpoint "/json/new" "创建新的浏览器页面"; then
    ((passed_tests++))
fi
((total_tests++))
echo ""

# WebSocket 连接测试
echo "======================================"
echo "测试 5: WebSocket 连接"
echo "======================================"
echo -e "${BLUE}测试:${NC} WebSocket 连接能力"

# 获取 WebSocket URL
ws_url=$(curl -s "${TEST_URL}/json/version" | grep -o '"webSocketDebuggerUrl":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$ws_url" ]; then
    echo -e "${GREEN}✓ 获取到 WebSocket URL:${NC}"
    echo "$ws_url"
    ((passed_tests++))
    
    # 检测 wscat 是否安装
    if command -v wscat &> /dev/null; then
        echo ""
        echo -e "${YELLOW}提示: 可以使用以下命令测试 WebSocket 连接:${NC}"
        echo "wscat -c \"$ws_url\""
    else
        echo ""
        echo -e "${YELLOW}提示: 安装 wscat 以测试 WebSocket 连接:${NC}"
        echo "npm install -g wscat"
        echo "wscat -c \"$ws_url\""
    fi
else
    echo -e "${RED}✗ 无法获取 WebSocket URL${NC}"
fi
((total_tests++))
echo ""

# 网络诊断
echo "======================================"
echo "网络诊断"
echo "======================================"

# 提取主机和端口
HOST=$(echo "$TEST_URL" | sed -E 's|^https?://||' | cut -d':' -f1 | cut -d'/' -f1)
PORT=$(echo "$TEST_URL" | sed -E 's|^https?://||' | cut -d':' -f2 | cut -d'/' -f1)

# 如果没有指定端口，根据协议设置默认端口
if [ "$HOST" = "$PORT" ]; then
    if [[ "$TEST_URL" == https* ]]; then
        PORT=443
    else
        PORT=80
    fi
fi

echo "主机: $HOST"
echo "端口: $PORT"
echo ""

# Ping 测试
echo -e "${BLUE}Ping 测试:${NC}"
if ping -c 3 "$HOST" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 主机可达${NC}"
    ping -c 3 "$HOST" | tail -2
else
    echo -e "${YELLOW}⚠ Ping 失败（可能主机禁用了 ICMP）${NC}"
fi
echo ""

# 端口测试
echo -e "${BLUE}端口测试:${NC}"
if command -v nc &> /dev/null; then
    if nc -zv "$HOST" "$PORT" 2>&1 | grep -q "succeeded"; then
        echo -e "${GREEN}✓ 端口 $PORT 开放${NC}"
    else
        echo -e "${RED}✗ 端口 $PORT 无法访问${NC}"
    fi
elif command -v telnet &> /dev/null; then
    if timeout 3 telnet "$HOST" "$PORT" 2>&1 | grep -q "Connected"; then
        echo -e "${GREEN}✓ 端口 $PORT 开放${NC}"
    else
        echo -e "${RED}✗ 端口 $PORT 无法访问${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 未安装 nc 或 telnet，无法测试端口${NC}"
fi
echo ""

# SSL/TLS 测试（如果是 HTTPS）
if [[ "$TEST_URL" == https* ]]; then
    echo -e "${BLUE}SSL/TLS 证书测试:${NC}"
    if command -v openssl &> /dev/null; then
        cert_info=$(echo | openssl s_client -connect "$HOST:$PORT" -servername "$HOST" 2>/dev/null | openssl x509 -noout -dates 2>/dev/null)
        if [ ! -z "$cert_info" ]; then
            echo -e "${GREEN}✓ SSL 证书有效${NC}"
            echo "$cert_info"
        else
            echo -e "${RED}✗ SSL 证书无效或无法获取${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ 未安装 openssl，无法测试证书${NC}"
    fi
    echo ""
fi

# 总结
echo "======================================"
echo "测试总结"
echo "======================================"
echo ""
echo -e "总测试数: ${BLUE}$total_tests${NC}"
echo -e "通过测试: ${GREEN}$passed_tests${NC}"
echo -e "失败测试: ${RED}$((total_tests - passed_tests))${NC}"
echo ""

if [ $passed_tests -eq $total_tests ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}🎉 所有测试通过！内网穿透配置成功！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${BLUE}下一步：${NC}"
    echo "1. 在您的代码中使用此地址连接 Chrome"
    echo "2. 配置安全措施（认证、IP 白名单等）"
    echo "3. 监控访问日志和性能指标"
    echo ""
    echo -e "${BLUE}使用示例：${NC}"
    echo ""
    echo "Go (chromedp):"
    echo "  allocCtx, _ := chromedp.NewRemoteAllocator(ctx, \"${TEST_URL%/}\")"
    echo ""
    echo "Python (Selenium):"
    echo "  chrome_options.add_experimental_option(\"debuggerAddress\", \"${HOST}:${PORT}\")"
    echo ""
    echo "Node.js (Puppeteer):"
    echo "  const browser = await puppeteer.connect({ browserWSEndpoint: '$ws_url' })"
    echo ""
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}❌ 部分测试失败，请检查配置${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
    echo -e "${YELLOW}故障排查建议：${NC}"
    echo ""
    echo "1. 检查 Chrome 容器是否运行:"
    echo "   docker-compose ps"
    echo "   docker logs snapup-chrome"
    echo ""
    echo "2. 检查内网穿透服务是否运行:"
    echo "   docker-compose logs frpc"
    echo "   # 或"
    echo "   docker-compose logs cloudflared"
    echo ""
    echo "3. 检查防火墙规则:"
    echo "   sudo ufw status"
    echo "   sudo iptables -L"
    echo ""
    echo "4. 检查 frp 服务端日志"
    echo ""
    echo "5. 验证 DNS 解析（如果使用域名）:"
    echo "   nslookup $HOST"
    echo "   dig $HOST"
    echo ""
    echo "详细文档: docs/TUNNEL_SETUP.md"
    echo ""
    exit 1
fi

