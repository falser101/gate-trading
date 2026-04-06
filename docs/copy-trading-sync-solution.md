# 跟单交易同步方案设计

## 一、背景与需求

**为什么需要这个方案：**
- Gate.io 官方 Go SDK 没有跟单交易 API
- 用户通过浏览器发现 Web API (`https://www.gate.com/apiw/v2/copy/follower/position`)
- Cookie 认证会过期，需要定时轮询 + 自动刷新机制
- 目标：为平台用户同步跟单数据（交易员持仓、盈亏等）

**核心需求：**
1. 存储用户的 Gate.io Cookie 认证信息（加密）
2. 定时轮询 Web API 获取跟单数据
3. Cookie 过期时自动刷新或通知用户
4. 将跟单数据同步到本地数据库供平台用户使用
5. 支持多个跟单账户

---

## 二、数据模型设计

### 2.1 CopyTradingAccount（跟单账户）

```go
type CopyTradingAccount struct {
    ID              uint           `gorm:"primarykey"`
    UserID          uint           `gorm:"index;not null"`  // 平台用户 ID
    GateUserId      string         `gorm:"size:50;not null"` // Gate.io 用户 ID（从 cookie 获取）
    CookieToken     string         `gorm:"size:2048;not null"` // 加密存储的 token cookie
    CookieCsrf      string         `gorm:"size:256"`        // csrftoken
    CookieExpiresAt *time.Time     `gorm:"index"`           // Cookie 过期时间
    Status          string         `gorm:"size:20;default:'active'"` // active/expired/error
    LastSyncedAt    *time.Time     // 最后同步时间
    SyncError       string         `gorm:"size:512"`        // 最近错误信息
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt
}
```

**字段说明：**
- `UserID`: 关联平台用户，支持多用户各自绑定跟单账户
- `CookieToken`: 存储 `token` cookie 值（JWT 格式，包含过期信息）
- `CookieCsrf`: 存储 `csrftoken` 用于请求验证
- `Status`: 
  - `active`: 正常同步中
  - `expired`: Cookie 过期，需要重新绑定
  - `error`: 同步出错（如 IP 限制、API 变更）

### 2.2 CopyTraderInfo（交易员信息）

```go
type CopyTraderInfo struct {
    ID                  uint        `gorm:"primarykey"`
    CopyTradingAccountID uint        `gorm:"index;not null"` // 关联跟单账户
    TraderId            int64       `gorm:"index;not null"`  // 交易员 ID
    TraderName          string      `gorm:"size:100"`        // 交易员昵称
    TraderAvatar        string      `gorm:"size:512"`        // 头像 URL
    PositionCount       int         `gorm:"default:0"`       // 持仓数量
    TotalPnl            string      `gorm:"size:32"`         // 总盈亏 (USDT)
    TotalRoi            string      `gorm:"size:32"`         // 总收益率
    WinRate             string      `gorm:"size:16"`         // 胜率
    LastPositionUpdate  *time.Time  // 最后持仓更新时间
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

### 2.3 CopyTradePosition（跟单持仓）

```go
type CopyTradePosition struct {
    ID                  uint        `gorm:"primarykey"`
    CopyTradingAccountID uint        `gorm:"index;not null"`
    TraderId            int64       `gorm:"index;not null"`
    PositionId          int64       `gorm:"uniqueIndex;not null"` // Gate.io Position ID
    Symbol              string      `gorm:"size:32;not null"`     // 交易对 BTC_USDT
    Side                string      `gorm:"size:10"`              // LONG/SHORT
    EntryPrice          string      `gorm:"size:32"`              // 开仓均价
    CurrentPrice        string      `gorm:"size:32"`              // 当前价格
    Size                string      `gorm:"size:32"`              // 持仓数量
    Leverage            string      `gorm:"size:16"`              // 杠杆倍数
    Pnl                 string      `gorm:"size:32"`              // 未实现盈亏
    Roi                 string      `gorm:"size:16"`              // 收益率
    OpenTime            *time.Time  // 开仓时间
    CloseTime           *time.Time  // 平仓时间
    Status              string      `gorm:"size:20;default:'open'"` // open/closed
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

### 2.4 CopyTradeHistory（跟单历史记录）

```go
type CopyTradeHistory struct {
    ID                  uint        `gorm:"primarykey"`
    CopyTradingAccountID uint        `gorm:"index;not null"`
    TraderId            int64       `gorm:"index;not null"`
    GatePositionId      int64       `gorm:"index"`
    Symbol              string      `gorm:"size:32"`
    Side                string      `gorm:"size:10"`
    EntryPrice          string      `gorm:"size:32"`
    ExitPrice           string      `gorm:"size:32"`
    Size                string      `gorm:"size:32"`
    Pnl                 string      `gorm:"size:32"`
    Roi                 string      `gorm:"size:16"`
    OpenTime            time.Time
    CloseTime           time.Time
    CreatedAt           time.Time
}
```

---

## 三、服务层设计

### 3.1 Cookie 管理服务

**文件：** `internal/service/copytrading/cookie_manager.go`

```go
type CookieManager struct {
    cryptoKey []byte // AES 加密密钥
}

// Cookie 结构
type GateCookies struct {
    Token     string    // token cookie (JWT)
    Csrf      string    // csrftoken
    Uid       string    // uid
    ExpiresAt time.Time // 从 token 中解析
}

// 方法
func (m *CookieManager) Encrypt(rawCookie string) (string, error)
func (m *CookieManager) Decrypt(encrypted string) (string, error)
func (m *CookieManager) ParseToken(token string) (*GateCookies, error)
func (m *CookieManager) IsExpired(cookies *GateCookies) bool
func (m *CookieManager) RefreshCookies(uid, token, csrf string) (*GateCookies, error)
```

**说明：**
- Cookie 使用 AES-256-GCM 加密存储
- Token 是 JWT 格式，可以解析出过期时间
- 刷新机制：尝试重新获取（需要用户重新登录 Gate.io）

### 3.2 HTTP 客户端服务

**文件：** `internal/service/copytrading/gate_web_client.go`

```go
type GateWebClient struct {
    httpClient *http.Client
    baseURL    string // https://www.gate.com
}

// API 响应结构
type FollowerPositionResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Positions []Position `json:"data"`
        Total     int        `json:"total"`
        Page      int        `json:"page"`
        PageSize  int        `json:"page_size"`
    } `json:"data"`
}

type Position struct {
    PositionId   int64   `json:"position_id"`
    TraderId     int64   `json:"trader_id"`
    TraderName   string  `json:"trader_name"`
    Symbol       string  `json:"market"`
    Side         string  `json:"side"` // 1=多，-1=空
    EntryPrice   float64 `json:"entry_price"`
    CurrentPrice float64 `json:"current_price"`
    Size         float64 `json:"size"`
    Leverage     int     `json:"leverage"`
    Pnl          float64 `json:"pnl"`
    Roi          float64 `json:"roi"`
    OpenTime     int64   `json:"open_time"`
}

// 方法
func (c *GateWebClient) GetFollowerPositions(cookies *GateCookies, page, pageSize int) (*FollowerPositionResponse, error)
func (c *GateWebClient) GetTraderInfo(cookies *GateCookies, traderId int64) (*TraderInfo, error)
func (c *GateWebClient) GetCopyTradingHistory(cookies *GateCookies, startDate, endDate time.Time) ([]Position, error)
```

**请求示例：**
```go
func (c *GateWebClient) GetFollowerPositions(cookies *GateCookies, page, pageSize int) (*FollowerPositionResponse, error) {
    url := fmt.Sprintf("https://www.gate.com/apiw/v2/copy/follower/position?trader_name=&market=&page=%d&page_size=%d", page, pageSize)
    
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Cookie", fmt.Sprintf("token=%s; csrftoken=%s; uid=%s", cookies.Token, cookies.Csrf, cookies.Uid))
    req.Header.Set("Csrftoken", cookies.Csrf)
    req.Header.Set("Referer", "https://www.gate.com/zh/copytrading/mine?mode=futures&type=copy")
    req.Header.Set("Sub-Website-Id", "0")
    req.Header.Set("User-Agent", "Mozilla/5.0...")
    
    resp, err := c.httpClient.Do(req)
    // ... 处理响应
}
```

### 3.3 同步服务

**文件：** `internal/service/copytrading/sync_service.go`

```go
type SyncService struct {
    repo           *CopyTradingRepository
    cookieManager  *CookieManager
    webClient      *GateWebClient
    syncInterval   time.Duration
    maxRetries     int
}

// 同步结果
type SyncResult struct {
    AccountID       uint
    SyncedAt        time.Time
    NewPositions    int
    UpdatedPositions int
    ClosedPositions int
    Error           error
}

// 方法
func (s *SyncService) SyncAccount(account *CopyTradingAccount) (*SyncResult, error)
func (s *SyncService) SyncAllAccounts() []*SyncResult
func (s *SyncService) CheckCookieExpiry(account *CopyTradingAccount) (bool, error)
func (s *SyncService) ProcessPosition(account *CopyTradingAccount, pos Position) error
```

**同步逻辑：**
1. 检查 Cookie 是否过期，过期则跳过
2. 分页获取所有跟单持仓
3. 对比本地数据库：
   - 新持仓：插入
   - 已存在：更新价格、盈亏
   - 本地有但 API 没有：标记为已平仓
4. 更新 `LastSyncedAt` 和时间戳

---

## 四、定时任务设计

### 4.1 Worker 实现

**文件：** `cmd/copytrading-worker/main.go`

```go
package main

import (
    "log"
    "time"
    "context"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/falser101/gate-trading/internal/config"
    "github.com/falser101/gate-trading/internal/database"
    "github.com/falser101/gate-trading/internal/service/copytrading"
)

func main() {
    // 加载配置
    cfg := config.Load()
    
    // 初始化数据库
    db := database.Init(cfg.DatabaseURL)
    
    // 初始化服务
    cookieManager := copytrading.NewCookieManager([]byte(cfg.EncryptionKey))
    webClient := copytrading.NewGateWebClient()
    syncService := copytrading.NewSyncService(db, cookieManager, webClient, 30*time.Second, 3)
    
    // 定时同步
    ticker := time.NewTicker(30 * time.Second) // 每 30 秒同步一次
    defer ticker.Stop()
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 优雅关闭
    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig
        cancel()
    }()
    
    log.Println("Copy trading sync worker started")
    
    for {
        select {
        case <-ticker.C:
            results := syncService.SyncAllAccounts()
            for _, r := range results {
                if r.Error != nil {
                    log.Printf("Sync account %d failed: %v", r.AccountID, r.Error)
                } else {
                    log.Printf("Sync account %d: +%d new, %d updated, %d closed", 
                        r.AccountID, r.NewPositions, r.UpdatedPositions, r.ClosedPositions)
                }
            }
        case <-ctx.Done():
            log.Println("Worker shutting down")
            return
        }
    }
}
```

### 4.2 Cookie 过期处理

**策略：**
1. **检测过期：** 每次同步前检查 Token JWT 的 `exp` 字段
2. **提前刷新：** 过期前 1 小时尝试刷新
3. **刷新失败：** 
   - 标记账户状态为 `expired`
   - 发送通知给用户（邮件/站内信）
   - 暂停同步直到用户重新绑定

**刷新机制选项：**

| 方案 | 优点 | 缺点 |
|------|------|------|
| 用户重新登录 Gate.io 获取新 Cookie | 可靠 | 需要用户手动操作 |
| 使用 refresh token 自动刷新 | 自动化 | Gate.io Web API 可能不支持 |
| 使用官方 API Key 替代 | 稳定、官方支持 | 无法获取跟单数据（SDK 无此功能） |

**推荐：** 方案 1 + 通知机制

---

## 五、API 接口设计

### 5.1 账户管理

**文件：** `cmd/api/main.go`

```go
// 绑定跟单账户
POST /api/copytrading/account
Body: {
    "cookie_token": "eyJ...",
    "cookie_csrf": "63304a56...",
    "uid": "49213049"
}

// 获取账户状态
GET /api/copytrading/account

// 更新 Cookie
PUT /api/copytrading/account
Body: {
    "cookie_token": "...",
    "cookie_csrf": "..."
}

// 删除账户
DELETE /api/copytrading/account

// 手动触发同步
POST /api/copytrading/account/sync

// 测试 Cookie 是否有效
POST /api/copytrading/account/test
```

### 5.2 数据查询

```go
// 获取跟单持仓列表
GET /api/copytrading/positions?status=open&page=1&page_size=20

// 获取单个交易员详情
GET /api/copytrading/trader/:id

// 获取跟单历史记录
GET /api/copytrading/history?start_date=2026-01-01&end_date=2026-04-06

// 获取统计数据
GET /api/copytrading/stats
Response: {
    "total_pnl": "1234.56",
    "total_roi": "12.34",
    "win_rate": "65.5",
    "position_count": 5,
    "trader_count": 3
}
```

---

## 六、前端设计（Flutter）

### 6.1 页面结构

**文件：** `lib/presentation/screens/copytrading/`

```
copytrading/
├── account_screen.dart       # 账户绑定/管理
├── positions_screen.dart     # 跟单持仓列表
├── traders_screen.dart       # 交易员列表
├── history_screen.dart       # 历史记录
└── widgets/
    ├── position_card.dart    # 持仓卡片
    ├── trader_card.dart      # 交易员卡片
    └── stats_panel.dart      # 统计面板
```

### 6.2 账户绑定页面

```
┌─────────────────────────────────────┐
│  绑定跟单账户                       │
├─────────────────────────────────────┤
│                                     │
│  请按照以下步骤获取 Cookie：          │
│                                     │
│  1. 在浏览器登录 Gate.io             │
│  2. 进入「跟单交易」→「我的跟单」    │
│  3. 按 F12 打开开发者工具            │
│  4. 找到 position 接口请求           │
│  5. 复制 Cookie 中的以下值:           │
│     - token                         │
│     - csrftoken                     │
│     - uid                           │
│                                     │
│  Token:  [___________________]      │
│  Csrf:   [__________]               │
│  UID:    [__________]               │
│                                     │
│  [取消]              [立即绑定]     │
└─────────────────────────────────────┘
```

### 6.3 持仓列表页面

```
┌─────────────────────────────────────┐
│  跟单持仓                    [刷新] │
├─────────────────────────────────────┤
│  总盈亏：+$1,234.56  (+12.3%) ↑    │
│  持仓数：5  |  交易员：3            │
├─────────────────────────────────────┤
│  ┌───────────────────────────────┐  │
│  │ BTC/USDT  交易员：CryptoKing  │  │
│  │ 多  10x  入场：65,000         │  │
│  │ 当前：67,500  盈亏：+$250     │  │
│  │ 收益率：+15.3%                │  │
│  └───────────────────────────────┘  │
│  ┌───────────────────────────────┐  │
│  │ ETH/USDT  交易员：EthMaster   │  │
│  │ 空  5x  入场：3,500           │  │
│  │ 当前：3,400  盈亏：+$100      │  │
│  │ 收益率：+8.5%                 │  │
│  └───────────────────────────────┘  │
│  ...                                │
└─────────────────────────────────────┘
```

### 6.4 Provider 状态管理

**文件：** `lib/presentation/providers/copytrading_provider.dart`

```dart
@riverpod
class CopytradingNotifier extends _$CopytradingNotifier {
  Future<CopytradingAccount?> getAccount();
  Future<void> bindAccount(String token, String csrf, String uid);
  Future<void> updateAccount(String token, String csrf);
  Future<void> deleteAccount();
  Future<void> syncNow();
  Future<List<CopyTradePosition>> getPositions({String? status});
  Future<List<CopyTrader>> getTraders();
  Future<CopytradingStats> getStats();
}
```

---

## 七、安全考虑

### 7.1 Cookie 加密存储

```go
// 使用 AES-256-GCM 加密
func (m *CookieManager) Encrypt(plaintext string) (string, error) {
    block, _ := aes.NewCipher(m.cryptoKey)
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    // 生成随机 nonce
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

### 7.2 环境配置

**.env 文件：**
```env
# 加密密钥（32 字节）
ENCRYPTION_KEY=your-32-byte-secret-key-here

# Worker 配置
COPYTRADING_SYNC_INTERVAL=30s
COPYTRADING_MAX_RETRIES=3
```

---

## 八、实现步骤

### Phase 1: 后端基础（预计 2-3 天）

1. **数据模型** (`internal/models/models.go`)
   - [ ] CopyTradingAccount
   - [ ] CopyTraderInfo
   - [ ] CopyTradePosition
   - [ ] CopyTradeHistory

2. **Repository 层** (`internal/repository/copytrading_repository.go`)
   - [ ] CRUD 操作
   - [ ] 查询优化（索引）

3. **服务层** (`internal/service/copytrading/`)
   - [ ] cookie_manager.go
   - [ ] gate_web_client.go
   - [ ] sync_service.go

4. **API 接口** (`cmd/api/main.go`)
   - [ ] 账户绑定/查询/更新/删除
   - [ ] 持仓查询
   - [ ] 统计数据

### Phase 2: Worker 实现（预计 1-2 天）

1. **定时同步 Worker** (`cmd/copytrading-worker/main.go`)
   - [ ] 定时任务循环
   - [ ] 错误处理
   - [ ] 优雅关闭

2. **Cookie 过期处理**
   - [ ] 过期检测
   - [ ] 通知机制

### Phase 3: 前端实现（预计 2-3 天）

1. **数据模型** (`lib/data/models/copytrading/`)
   - [ ] account_model.dart
   - [ ] position_model.dart
   - [ ] trader_model.dart

2. **Repository** (`lib/data/repositories/copytrading_repository.dart`)
   - [ ] API 调用封装

3. **Provider** (`lib/presentation/providers/copytrading_provider.dart`)
   - [ ] 状态管理

4. **页面** (`lib/presentation/screens/copytrading/`)
   - [ ] account_screen.dart
   - [ ] positions_screen.dart
   - [ ] traders_screen.dart
   - [ ] history_screen.dart

### Phase 4: 测试与优化（预计 1 天）

1. **后端测试**
   - [ ] API 接口测试
   - [ ] 同步逻辑测试
   - [ ] Cookie 过期测试

2. **前端测试**
   - [ ] 页面 UI 测试
   - [ ] 状态管理测试

3. **性能优化**
   - [ ] 数据库索引
   - [ ] 分页查询
   - [ ] 缓存策略

---

## 九、风险与注意事项

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Cookie 过期 | 同步中断 | 提前通知用户重新绑定 |
| Gate.io Web API 变更 | 功能失效 | 监控 API 响应，快速适配 |
| 频率限制 | 请求被限 | 合理设置同步间隔（30s+） |
| Cookie 泄露 | 安全风险 | AES 加密存储，限制访问权限 |
| IP 限制 | 无法访问 | 考虑使用代理或服务器 IP 白名单 |

---

## 十、扩展方向

1. **自动跟单：** 对接官方 API，实现自动跟随交易员下单
2. **交易员排行榜：** 抓取公开交易员数据，展示排名
3. **策略分析：** 分析交易员历史表现，给出推荐
4. **多账户管理：** 支持绑定多个 Gate.io 账户
5. **推送通知：** 重要变动（大额盈亏、平仓）推送给用户

---

## 关键文件路径总结

| 类型 | 文件路径 |
|------|----------|
| 数据模型 | `internal/models/models.go` |
| Repository | `internal/repository/copytrading_repository.go` |
| Cookie 管理 | `internal/service/copytrading/cookie_manager.go` |
| HTTP 客户端 | `internal/service/copytrading/gate_web_client.go` |
| 同步服务 | `internal/service/copytrading/sync_service.go` |
| Worker | `cmd/copytrading-worker/main.go` |
| API 路由 | `cmd/api/main.go` |
| Flutter 模型 | `flutter_app/lib/data/models/copytrading/` |
| Flutter Provider | `flutter_app/lib/presentation/providers/copytrading_provider.dart` |
| Flutter 页面 | `flutter_app/lib/presentation/screens/copytrading/` |
