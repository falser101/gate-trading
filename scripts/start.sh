#!/bin/bash
# Copy Trading 完整启动脚本
# 用法：./scripts/start.sh [api|worker|both|test|init]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}=== $1 ===${NC}"
}

print_error() {
    echo -e "${RED}错误：$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}提示：$1${NC}"
}

show_help() {
    echo "Copy Trading 启动脚本"
    echo ""
    echo "用法：$0 [command]"
    echo ""
    echo "命令:"
    echo "  api       启动 API 服务"
    echo "  worker    启动同步 Worker"
    echo "  both      同时启动 API 和 Worker"
    echo "  test      测试 Cookie"
    echo "  init      初始化 Cookie"
    echo "  sync      手动同步一次"
    echo "  status    查看 Cookie 状态"
    echo "  help      显示帮助"
    echo ""
    echo "示例:"
    echo "  $0 api      # 启动 API"
    echo "  $0 both     # 同时启动 API 和 Worker"
    echo "  $0 test     # 测试 Cookie 是否有效"
}

init_cookie() {
    print_info "初始化 Copy Trading Cookie"

    # Cookie 信息
    GATE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NzU0NDEzNzQsImlwIjoiT0Y3eTQ3TXQ2cXhpakJSYmN0MXRYNDdzZ2pTSUtEV2VBNk92SDZNT0M5SzdVVHF0eS9KaS9sOEMiLCJpcFJlc3RyaWN0IjoiYUdlcVJNYzJsQThwdjdGUWx0UEs2cFhvRFEyeDFLWnQreHNiK0ZZPSIsImRldmljZVR5cGUiOiI5dkNPTGNSSUs0OCtFK3JNWWIvRHN2b1A4N2ZteEtlQWpoSHA3aE09IiwiZGV2aWNlSWQiOiI3THIwOG5jWklqWlRRTXFnMnRybXRIeUFQbzV4WHFkVVUwVmt1QT05IiwidWlkIjoiZ0lUaGY1U3E5M2lsMVhsMll6YjNDWEV4c253M2Y3M2dweDVYVmxTaVhpSFFCNjU1Iiwid2Vic2l0ZUlEIjoiaXhadmhwd1g0cXJzWFNhUW1SZG0xV25yazR3bklGUlBjVExLRi93PSJ9.ys7rQyESjz1IMCm3vBWOR9dMmdeQdbdAKIAWb08xmNE"
    GATE_CSRFTOKEN="63304a56794365546b50794d49593343324d6b626e535649574371336a7a4f664f497631476f74726b4d6e6d62505534526175384557792b7945746c70456a34"
    GATE_UID="49213049"

    go run ./cmd/copytrading-sync-worker/main.go init "$GATE_TOKEN" "$GATE_CSRFTOKEN" "$GATE_UID"
}

test_cookie() {
    print_info "测试 Copy Trading Cookie"
    go run ./cmd/copytrading-sync-worker/main.go test
}

run_api() {
    print_info "启动 API 服务"
    print_warning "按 Ctrl+C 停止服务"
    go run ./cmd/api/
}

run_worker() {
    print_info "启动同步 Worker"
    print_warning "按 Ctrl+C 停止服务"
    go run ./cmd/copytrading-sync-worker/main.go
}

run_both() {
    print_info "同时启动 API 和 Worker"
    print_warning "按 Ctrl+C 停止所有服务"

    # 启动 Worker（后台）
    go run ./cmd/copytrading-sync-worker/main.go &
    WORKER_PID=$!

    # 等待 2 秒
    sleep 2

    # 启动 API（前台）
    trap "kill $WORKER_PID 2>/dev/null" EXIT
    go run ./cmd/api/
}

manual_sync() {
    print_info "手动同步交易员数据"
    go run ./cmd/copytrading-sync-worker/main.go sync
}

check_status() {
    print_info "Cookie 状态"

    # 检查 .env 文件
    if [ -f ".env" ]; then
        source .env 2>/dev/null || true
        if [ -n "$ENCRYPTION_KEY" ]; then
            echo "  加密密钥：已配置"
        else
            echo "  加密密钥：未配置"
        fi
    fi

    # 检查 Worker 是否运行
    if pgrep -f "copytrading-sync-worker" > /dev/null; then
        echo "  Worker 状态：运行中"
    else
        echo "  Worker 状态：未运行"
    fi

    # 检查 API 是否运行
    if pgrep -f "cmd/api" > /dev/null; then
        echo "  API 状态：运行中"
    else
        echo "  API 状态：未运行"
    fi
}

# 主逻辑
case "${1:-help}" in
    api)
        run_api
        ;;
    worker)
        run_worker
        ;;
    both)
        run_both
        ;;
    test)
        test_cookie
        ;;
    init)
        init_cookie
        ;;
    sync)
        manual_sync
        ;;
    status)
        check_status
        ;;
    help|-h|--help)
        show_help
        ;;
    *)
        print_error "未知命令：$1"
        show_help
        exit 1
        ;;
esac
