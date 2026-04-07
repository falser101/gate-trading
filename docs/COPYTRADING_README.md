# Copy Trading 交易员数据同步

## 快速开始

### 1. 配置环境变量

复制 `.env.example` 为 `.env` 并配置：

```bash
# 同步间隔（分钟），默认 30 分钟
COPYTRADING_SYNC_INTERVAL=30

# 加密密钥（32 字节，用于加密存储 Cookie）
ENCRYPTION_KEY=your-32-byte-encryption-key-here!!

# User-Agent（建议使用真实浏览器 UA）
COPYTRADING_USER_AGENT=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36...
```

### 2. 初始化 Cookie

管理员需要先从 Gate.io 浏览器复制 Cookie，然后初始化：

```bash
# 运行初始化命令
go run cmd/copytrading-sync-worker/main.go init <token> <csrftoken> <uid>
```

**如何获取 Cookie：**
1. 浏览器登录 https://www.gate.com
2. 按 F12 打开开发者工具 → Network 标签
3. 访问 https://www.gate.com/zh/copytrading
4. 找到 `leader/list` 请求
5. 复制 Request Headers 中的：
   - `token` (JWT)
   - `csrftoken`
   - `uid`

### 3. 测试 Cookie

```bash
# 测试 Cookie 是否有效
go run cmd/copytrading-sync-worker/main.go test
```

### 4. 启动同步 Worker

```bash
# 后台运行
go run cmd/copytrading-sync-worker/main.go &

# 或编译后运行
go build -o copytrading-sync-worker ./cmd/copytrading-sync-worker/
./copytrading-sync-worker
```

### 5. 启动 API 服务

```bash
go run ./cmd/api/
```

## API 接口

### 公开接口（需要用户认证）

```bash
# 获取交易员列表
GET /api/copytrading/traders?page=1&page_size=20&order_by=follow_profit&cycle=month

# 获取交易员详情
GET /api/copytrading/traders/:id

# 获取交易员历史统计
GET /api/copytrading/traders/:id/stats?start_date=2026-01-01&end_date=2026-04-07
```

### 管理员接口

```bash
# 保存 Cookie
POST /api/admin/copytrading/cookie
{
  "token": "eyJhbGc...",
  "csrftoken": "63304...",
  "uid": "49213049"
}

# 获取 Cookie 状态
GET /api/admin/copytrading/cookie

# 删除 Cookie
DELETE /api/admin/copytrading/cookie

# 手动触发同步
POST /api/admin/copytrading/sync

# 获取管理员通知
GET /api/admin/notifications
```

## Worker 命令

```bash
# 运行定时同步
copytrading-sync-worker

# 测试 Cookie
copytrading-sync-worker test

# 初始化 Cookie
copytrading-sync-worker init <token> <csrftoken> <uid>

# 手动同步
copytrading-sync-worker sync

# 显示帮助
copytrading-sync-worker help
```

## 数据模型

### PlatformCookie
- 存储 Gate.io Cookie（加密）
- 记录过期时间和同步状态

### CopyTrader
- 交易员基本信息
- 统计数据（PnL、ROI、胜率等）
- 风格标签

### CopyTraderDailyStats
- 每日快照数据
- 用于绘制历史收益曲线

### AdminNotification
- Cookie 过期提醒
- 系统通知

## 注意事项

1. **Cookie 有效期**：Token JWT 通常 7-30 天过期，过期前 7 天会收到通知
2. **同步频率**：建议 30 分钟以上，避免触发 Gate.io 风控
3. **加密存储**：Cookie 使用 AES-256-GCM 加密，确保安全
4. **IP 限制**：如遇到访问问题，考虑使用住宅 IP 代理
