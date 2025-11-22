#!/bin/bash

# SnapUp MCP Server 测试脚本
# 此脚本可用于测试 MCP Server 的基本功能

echo "启动 SnapUp MCP Server..."
echo "请确保已经构建了 snapup 可执行文件"
echo ""

# 进入项目根目录
cd "$(dirname "$0")/.."

# 确保 screenshots 目录存在
mkdir -p ./screenshots

# 启动 MCP Server
echo "运行 MCP Server（按 Ctrl+C 退出）"
echo "---"
echo ""

# 测试初始化请求
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}' | ./snapup -mode=mcp

echo ""
echo "MCP Server 已退出"

