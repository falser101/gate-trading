# 跟单交易 - 交易员数据同步方案

## 一、背景与目标

Gate.io 官方 API V4 没有提供查询交易员信息的接口，但 Web 端提供了内部 API：

```
GET https://www.gate.com/apiw/v2/copy/leader/list
```

该接口需要登录 Cookie 认证。本方案实现：
1. 平台管理员绑定 Gate.io Cookie
2. 后台定时同步交易员数据到本地数据库
3. 所有注册用户可从本地数据库查看交易员列表

---

## 二、API 分析

### 2.1 请求示例

```bash
curl 'https://www.gate.com/apiw/v2/copy/leader/list?sub_website_id=0&order_by=follow_profit&sort_by=desc&cycle=month&status=running&page=1&page_size=10' \
  -H 'cookie: token=eyJhbG...; csrftoken=63304...; uid=49213049' \
  -H 'csrftoken: 63304a56794365546b50794d49593343...' \
  -H 'sub_website_id: 0'
```

### 2.2 关键参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `order_by` | 排序字段：`follow_profit`/`roi`/`follower_count` | `follow_profit` |
| `sort_by` | 排序方向：`asc`/`desc` | `desc` |
| `cycle` | 统计周期：`day`/`week`/`month`/`all` | `month` |
| `status` | 状态：`running`/`stopped` | `running` |
| `page` | 页码 | `1` |
| `page_size` | 每页数量（最大 100） | `10` |

### 2.3 必需的 Cookie

| Cookie | 说明 | 获取方式 |
|--------|------|----------|
| `token` | JWT 认证令牌 | 浏览器登录后 F12 复制 |
| `csrftoken` | CSRF 令牌 | 同上 |
| `uid` | 用户 ID | 同上 |

---

## 三、数据模型设计

### 3.1 PlatformCookie（平台级 Cookie 配置）

```go
type PlatformCookie struct {
    ID          uint           `gorm:"primarykey"`
    Token       string         `gorm:"size:2048;not null"`  // 加密存储
    CsrfToken   string         `gorm:"size:256;not null"`   // 加密存储
    Uid         string         `gorm:"size:50;not null"`
    ExpiresAt   *time.Time     `gorm:"index"`               // JWT 过期时间
    Status      string         `gorm:"size:20;default:'active'"` // active/expired
    LastSyncedAt *time.Time
    ErrorMsg    string         `gorm:"size:512"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**说明：** 全平台只有一条有效 Cookie 记录

### 3.2 CopyTrader（交易员信息）

```go
type CopyTrader struct {
    ID              uint        `gorm:"primarykey"`
    TraderID        string      `gorm:"size:100;not null;uniqueIndex"`
    TraderName      string      `gorm:"size:100"`
    Avatar          string      `gorm:"size:512"`
    Exchange        string      `gorm:"size:20;default:'gate'"`
    Status          string      `gorm:"size:20;default:'running'"`
    
    // 统计数据
    Cycle           string      `gorm:"size:20"`
    TotalPnl        string      `gorm:"size:32"`
    TotalRoi        string      `gorm:"size:32"`
    FollowProfit    string      `gorm:"size:32"`
    FollowRoi       string      `gorm:"size:32"`
    WinRate         string      `gorm:"size:16"`
    FollowerCount   int         `gorm:"default:0"`
    PositionCount   int         `gorm:"default:0"`
    MaxDrawdown     string      `gorm:"size:32"`
    AvgLeverage     string      `gorm:"size:16"`
    
    // 属性
    IsCurated       bool        `gorm:"default:false"`
    IsPrivate       bool        `gorm:"default:false"`
    StyleLabels     string      `gorm:"type:jsonb"`  // JSON 数组
    
    // 同步
    LastSyncedAt    *time.Time  `gorm:"index"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 3.3 CopyTraderDailyStats（每日快照）

```go
type CopyTraderDailyStats struct {
    ID              uint      `gorm:"primarykey"`
    TraderID        string    `gorm:"size:100;index"`
    Date            string    `gorm:"size:10;index"`  // YYYY-MM-DD
    TotalPnl        string    `gorm:"size:32"`
    TotalRoi        string    `gorm:"size:32"`
    FollowProfit    string    `gorm:"size:32"`
    FollowerCount   int       `gorm:"default:0"`
    CreatedAt       time.Time
}
```

---

## 四、服务层设计

### 4.1 Cookie 管理

```go
type CookieManager struct {
    cryptoKey []byte
}

func (m *CookieManager) Encrypt(s string) (string, error)
func (m *CookieManager) Decrypt(s string) (string, error)
func (m *CookieManager) ParseToken(token string) (*time.Time, error) // 返回过期时间
func (m *CookieManager) ValidateCookie(cookie *GateCookie) (bool, error)
```

### 4.2 Gate Web 客户端

```go
type GateWebClient struct {
    httpClient *http.Client
}

func (c *GateWebClient) GetTraderList(cookie *GateCookie, params *ListParams) (*TraderListResponse, error)
```

### 4.3 同步服务

```go
type SyncService struct {
    repo      *CopyTradingRepository
    cookieMgr *CookieManager
    client    *GateWebClient
}

// 同步所有交易员数据（分页获取）
func (s *SyncService) SyncAll() (*SyncResult, error)

// 检查 Cookie 是否即将过期（< 7 天）
func (s *SyncService) CheckCookieExpiry() (bool, *time.Time, error)
```

---

## 五、定时任务设计

### 5.1 Worker 配置

```go
// cmd/copytrading-sync-worker/main.go

syncInterval := 30 * time.Minute  // 每 30 分钟同步一次
ticker := time.NewTicker(syncInterval)

for {
    select {
    case <-ticker.C:
        syncSvc.SyncAll()
    }
}
```

### 5.2 Cookie 过期检查

```go
// 每次同步前检查
cookie, err := repo.GetActiveCookie()
if err != nil {
    log.Println("No active cookie found")
    return
}

// 检查是否过期
if cookie.ExpiresAt.Before(time.Now()) {
    repo.UpdateStatus(cookie.ID, "expired")
    notifyAdmin("Cookie 已过期，请重新绑定")
    return
}

// 检查是否即将过期（< 7 天）
if cookie.ExpiresAt.Before(time.Now().Add(7 * 24 * time.Hour)) {
    notifyAdmin("Cookie 即将过期，请尽快重新绑定")
}
```

### 5.3 管理员通知

```go
// 简单实现：记录到数据库通知表
func notifyAdmin(message string) {
    repo.CreateAdminNotification(&AdminNotification{
        Type:    "cookie_expiry",
        Message: message,
        IsRead:  false,
    })
}

// 后续可扩展：邮件、Telegram、钉钉等
```

---

## 六、API 接口设计

### 6.1 管理员接口

```go
// 绑定/更新 Cookie
POST /api/admin/copytrading/cookie
Body: {
    "token": "eyJhbGc...",
    "csrftoken": "63304...",
    "uid": "49213049"
}

// 获取 Cookie 状态
GET /api/admin/copytrading/cookie
Response: {
    "status": "active",
    "expires_at": "2026-05-07T12:00:00Z",
    "last_synced_at": "2026-04-07T10:00:00Z"
}

// 手动触发同步
POST /api/admin/copytrading/sync

// 获取管理员通知列表
GET /api/admin/notifications
```

### 6.2 公开接口（所有用户可访问）

```go
// 获取交易员列表
GET /api/copytrading/traders?page=1&page_size=20&order_by=follow_profit&cycle=month
Response: {
    "data": [...],
    "total": 150,
    "page": 1,
    "page_size": 20
}

// 获取交易员详情
GET /api/copytrading/traders/:id

// 获取交易员历史统计
GET /api/copytrading/traders/:id/stats?start_date=2026-01-01&end_date=2026-04-07
```

---

## 七、管理员操作指南

### 7.1 获取并绑定 Cookie

1. 浏览器登录 https://www.gate.com
2. 按 F12 打开开发者工具 → Network 标签
3. 访问 https://www.gate.com/zh/copytrading
4. 找到 `leader/list` 请求
5. 复制 Request Headers 中的：
   - `cookie` 中的 `token=xxx; csrftoken=xxx; uid=xxx`
   - `csrftoken` header 值
6. 在管理平台粘贴保存

### 7.2 Cookie 有效期

- Token JWT 通常有效期 **7-30 天**
- 过期前 7 天系统会通知管理员
- 过期后需重新按上述步骤获取

---

## 八、实现步骤

### Phase 1: 后端基础（1-2 天）
- [ ] 数据模型：`PlatformCookie`, `CopyTrader`, `CopyTraderDailyStats`
- [ ] Repository 层
- [ ] CookieManager（加密/解密/JWT 解析）
- [ ] GateWebClient

### Phase 2: 同步与定时任务（1 天）
- [ ] SyncService
- [ ] Worker 实现
- [ ] Cookie 过期检查
- [ ] 管理员通知

### Phase 3: API 接口（0.5 天）
- [ ] 管理员接口（Cookie 管理）
- [ ] 公开接口（交易员查询）

### Phase 4: 前端（1-2 天）
- [ ] 管理员 Cookie 配置页
- [ ] 交易员列表页（Flutter）
- [ ] 交易员详情页

---

## 九、关键文件清单

| 文件 | 路径 |
|------|------|
| 数据模型 | `internal/models/models.go` |
| Repository | `internal/repository/copytrading_repository.go` |
| Cookie 管理 | `internal/service/copytrading/cookie_manager.go` |
| HTTP 客户端 | `internal/service/copytrading/gate_web_client.go` |
| 同步服务 | `internal/service/copytrading/sync_service.go` |
| Worker | `cmd/copytrading-sync-worker/main.go` |
| API 路由 | `cmd/api/main.go` |

---

## 十、环境变量

```env
# 加密密钥（32 字节）
ENCRYPTION_KEY=your-32-byte-secret-key-here

# 同步间隔
COPYTRADING_SYNC_INTERVAL=30m

# 数据库
DATABASE_URL=postgres://user:pass@localhost:5432/gate_trading
```
