#!/bin/bash
# Copy Trading Cookie 初始化脚本
# 自动从提供的 curl 命令中解析 Cookie

set -e

echo "=== Copy Trading Cookie 初始化 ==="
echo ""

# Cookie 信息（从 curl 命令解析）
GATE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NzU0NDEzNzQsImlwIjoiT0Y3eTQ3TXQ2cXhpakJSYmN0MXRYNDdzZ2pTSUtEV2VBNk92SDZNT0M5SzdVVHF0eS9KaS9sOEMiLCJpcFJlc3RyaWN0IjoiYUdlcVJNYzJsQThwdjdGUWx0UEs2cFhvRFEyeDFLWnQreHNiK0ZZPSIsImRldmljZVR5cGUiOiI5dkNPTGNSSUs0OCtFK3JNWWIvRHN2b1A4N2ZteEtlQWpoSHA3aE09IiwiZGV2aWNlSWQiOiI3THIwOG5jWklqWlRRTXFnMnRybXRIeUFQbzV4WHFkVVUwVmt1QT05IiwidWlkIjoiZ0lUaGY1U3E5M2lsMVhsMll6YjNDWEV4c253M2Y3M2dweDVYVmxTaVhpSFFCNjU1Iiwid2Vic2l0ZUlEIjoiaXhadmhwd1g0cXJzWFNhUW1SZG0xV25yazR3bklGUlBjVExLRi93PSJ9.ys7rQyESjz1IMCm3vBWOR9dMmdeQdbdAKIAWb08xmNE"
GATE_CSRFTOKEN="63304a56794365546b50794d49593343324d6b626e535649574371336a7a4f664f497631476f74726b4d6e6d62505534526175384557792b7945746c70456a34"
GATE_UID="49213049"

echo "使用以下 Cookie 进行初始化："
echo "  UID: $GATE_UID"
echo "  Token: ${GATE_TOKEN:0:50}..."
echo "  CsrfToken: ${GATE_CSRFTOKEN:0:30}..."
echo ""

# 执行初始化
echo "正在初始化 Cookie..."
go run ./cmd/copytrading-sync-worker/main.go init "$GATE_TOKEN" "$GATE_CSRFTOKEN" "$GATE_UID"

echo ""
echo "=== 初始化完成 ==="
echo ""
echo "接下来可以运行："
echo "  go run ./cmd/copytrading-sync-worker/main.go test    # 测试 Cookie"
echo "  go run ./cmd/copytrading-sync-worker/main.go         # 启动同步 Worker"
echo "  go run ./cmd/api/                                    # 启动 API 服务"
