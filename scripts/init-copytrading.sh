#!/bin/bash
# Copy Trading Cookie 初始化脚本
# 用法：./scripts/init-copytrading.sh

# 从 .env.local 读取 Cookie
set -a
source .env.local 2>/dev/null || true
set +a

echo "=== Copy Trading Cookie 初始化 ==="
echo ""

# 检查是否存在 Cookie 配置
if [ -z "$GATE_COOKIE_TOKEN" ] || [ -z "$GATE_COOKIE_CSRFTOKEN" ] || [ -z "$GATE_COOKIE_UID" ]; then
    echo "错误：请在 .env.local 文件中配置以下变量："
    echo "  GATE_COOKIE_TOKEN=..."
    echo "  GATE_COOKIE_CSRFTOKEN=..."
    echo "  GATE_COOKIE_UID=..."
    echo ""
    echo "或者直接在命令行运行："
    echo "go run ./cmd/copytrading-sync-worker/main.go init <token> <csrftoken> <uid>"
    exit 1
fi

echo "使用以下 Cookie 进行初始化："
echo "  UID: $GATE_COOKIE_UID"
echo "  Token: ${GATE_COOKIE_TOKEN:0:50}..."
echo "  CsrfToken: ${GATE_COOKIE_CSRFTOKEN:0:30}..."
echo ""

# 执行初始化
go run ./cmd/copytrading-sync-worker/main.go init "$GATE_COOKIE_TOKEN" "$GATE_COOKIE_CSRFTOKEN" "$GATE_COOKIE_UID"
