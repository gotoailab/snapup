#!/bin/bash

# SnapUp MCP Docker 管理脚本
# 用于快速启动、停止和管理 Docker 中的 MCP Server

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
COMPOSE_FILE="docker-compose.yml"
COMPOSE_FILE_CN="docker-compose.cn.yml"
MCP_SERVICE="snapup-mcp"
CHROME_SERVICE="chrome"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
SnapUp MCP Docker 管理脚本

用法:
    $0 <command> [options]

命令:
    start       启动 MCP 服务
    stop        停止 MCP 服务
    restart     重启 MCP 服务
    logs        查看 MCP 服务日志
    status      查看服务状态
    build       重新构建 MCP 服务镜像
    clean       停止并删除所有容器
    help        显示此帮助信息

选项:
    --cn        使用中国版配置（国内镜像源）
    -f          实时跟踪日志（仅用于 logs 命令）

示例:
    $0 start            # 启动 MCP 服务
    $0 start --cn       # 使用国内镜像源启动
    $0 logs -f          # 实时查看日志
    $0 build --cn       # 使用国内镜像源重新构建
    $0 status           # 查看服务状态

EOF
}

# 检测是否使用中国版配置
detect_compose_file() {
    USE_CN=false
    for arg in "$@"; do
        if [ "$arg" = "--cn" ]; then
            USE_CN=true
            break
        fi
    done

    if [ "$USE_CN" = true ]; then
        COMPOSE_CMD="docker-compose -f $COMPOSE_FILE_CN"
        print_info "使用中国版配置文件: $COMPOSE_FILE_CN"
    else
        COMPOSE_CMD="docker-compose -f $COMPOSE_FILE"
        print_info "使用标准配置文件: $COMPOSE_FILE"
    fi
}

# 检查 Docker 是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker 未运行，请先启动 Docker"
        exit 1
    fi
}

# 检查 docker-compose 是否安装
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose 未安装，请先安装 docker-compose"
        exit 1
    fi
}

# 启动服务
start_service() {
    print_info "正在启动 MCP 服务..."
    
    # 检查 Chrome 服务是否运行
    if ! $COMPOSE_CMD ps $CHROME_SERVICE | grep -q "Up"; then
        print_info "Chrome 服务未运行，正在启动 Chrome..."
        $COMPOSE_CMD up -d $CHROME_SERVICE
        print_info "等待 Chrome 启动（10秒）..."
        sleep 10
    fi
    
    # 启动 MCP 服务
    $COMPOSE_CMD --profile mcp up -d $MCP_SERVICE
    
    print_success "MCP 服务已启动"
    print_info "运行 '$0 logs -f' 查看实时日志"
}

# 停止服务
stop_service() {
    print_info "正在停止 MCP 服务..."
    $COMPOSE_CMD stop $MCP_SERVICE
    print_success "MCP 服务已停止"
}

# 重启服务
restart_service() {
    print_info "正在重启 MCP 服务..."
    $COMPOSE_CMD restart $MCP_SERVICE
    print_success "MCP 服务已重启"
}

# 查看日志
show_logs() {
    FOLLOW=""
    for arg in "$@"; do
        if [ "$arg" = "-f" ] || [ "$arg" = "--follow" ]; then
            FOLLOW="-f"
            break
        fi
    done
    
    if [ -n "$FOLLOW" ]; then
        print_info "实时查看 MCP 服务日志（Ctrl+C 退出）..."
        $COMPOSE_CMD logs $FOLLOW $MCP_SERVICE
    else
        print_info "查看 MCP 服务日志..."
        $COMPOSE_CMD logs --tail=100 $MCP_SERVICE
    fi
}

# 查看状态
show_status() {
    print_info "服务状态:"
    echo ""
    $COMPOSE_CMD ps
    echo ""
    
    # 检查 MCP 容器是否运行
    if $COMPOSE_CMD ps $MCP_SERVICE | grep -q "Up"; then
        print_success "MCP 服务正在运行"
    else
        print_warning "MCP 服务未运行"
    fi
    
    # 检查 Chrome 容器是否运行
    if $COMPOSE_CMD ps $CHROME_SERVICE | grep -q "Up"; then
        print_success "Chrome 服务正在运行"
    else
        print_warning "Chrome 服务未运行"
    fi
}

# 重新构建
rebuild_service() {
    print_info "正在重新构建 MCP 服务镜像..."
    $COMPOSE_CMD build --no-cache $MCP_SERVICE
    print_success "镜像构建完成"
    print_info "运行 '$0 start' 启动服务"
}

# 清理
clean_service() {
    print_warning "这将停止并删除所有容器，确定吗？(y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_info "正在清理..."
        $COMPOSE_CMD down
        print_success "清理完成"
    else
        print_info "操作已取消"
    fi
}

# 主函数
main() {
    # 检查基本依赖
    check_docker
    check_docker_compose
    
    # 检测配置文件
    detect_compose_file "$@"
    
    # 处理命令
    COMMAND="${1:-help}"
    
    case "$COMMAND" in
        start)
            start_service
            ;;
        stop)
            stop_service
            ;;
        restart)
            restart_service
            ;;
        logs)
            show_logs "$@"
            ;;
        status)
            show_status
            ;;
        build)
            rebuild_service
            ;;
        clean)
            clean_service
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $COMMAND"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"

