package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Email         string         `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash  string         `gorm:"size:255" json:"-"`
	GateAPIKey    string         `gorm:"size:255" json:"-"`
	GateAPISecret string         `gorm:"size:255" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Strategies    []Strategy     `gorm:"foreignKey:UserID" json:"strategies,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Strategy struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	Name        string         `gorm:"size:100" json:"name"`
	Type        string         `gorm:"size:20;index" json:"type"` // grid, dca, combo
	Symbol      string         `gorm:"size:20;index" json:"symbol"` // BTC_USDT
	Config      ConfigJSON     `gorm:"type:jsonb" json:"config"`
	Status      string         `gorm:"size:20;default:stopped" json:"status"` // running, stopped, paused
	RunState    RunStateJSON   `gorm:"type:jsonb" json:"run_state,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Orders      []Order        `gorm:"foreignKey:StrategyID" json:"orders,omitempty"`
}

func (Strategy) TableName() string {
	return "strategies"
}

// 策略配置 - 网格策略
type GridConfig struct {
	LowerPrice   string `json:"lower_price"` // 价格下限
	UpperPrice   string `json:"upper_price"` // 价格上限
	GridCount    int    `json:"grid_count"`  // 网格数量
	InvestAmount string `json:"invest_amount"` // 投资金额
	ProfitRate   string `json:"profit_rate"`   // 单个网格利润率
}

// 策略配置 - DCA 策略
type DCAConfig struct {
	InvestAmount string `json:"invest_amount"` // 每次投资金额
	Interval     int    `json:"interval"`      // 执行间隔 (分钟)
	TargetPrice  string `json:"target_price"`  // 目标价格 (可选)
	MaxBuyTimes  int    `json:"max_buy_times"` // 最大买入次数
}

// 策略配置 - 组合策略
type ComboConfig struct {
	GridConfig
	DCAConfig
}

// 策略运行状态
type RunState struct {
	TotalProfit    string   `json:"total_profit"`    // 总盈亏
	TotalBuyTimes  int      `json:"total_buy_times"` // 总买入次数
	TotalSellTimes int      `json:"total_sell_times"` // 总卖出次数
	LastRunTime    *time.Time `json:"last_run_time"` // 最后运行时间
	ActiveOrders   []string `json:"active_orders"`   // 活跃订单 ID 列表
}

type ConfigJSON map[string]interface{}
type RunStateJSON map[string]interface{}

type Order struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	StrategyID   uint           `gorm:"index" json:"strategy_id"`
	UserID       uint           `gorm:"index" json:"user_id"`
	GateOrderID  string         `gorm:"size:100;index" json:"gate_order_id"`
	Symbol       string         `gorm:"size:20" json:"symbol"`
	Side         string         `gorm:"size:10" json:"side"` // buy, sell
	Type         string         `gorm:"size:10" json:"type"` // limit, market
	Price        string         `gorm:"size:30" json:"price"`
	Amount       string         `gorm:"size:30" json:"amount"`
	FilledAmount string         `gorm:"size:30" json:"filled_amount"`
	Status       string         `gorm:"size:20" json:"status"` // open, filled, cancelled
	Fee          string         `gorm:"size:30" json:"fee"`
	FeeCurrency  string         `gorm:"size:20" json:"fee_currency"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Order) TableName() string {
	return "orders"
}

type Trade struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	OrderID     uint           `gorm:"index" json:"order_id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	GateTradeID int64          `gorm:"index" json:"gate_trade_id"`
	Symbol      string         `gorm:"size:20" json:"symbol"`
	Side        string         `gorm:"size:10" json:"side"`
	Price       string         `gorm:"size:30" json:"price"`
	Amount      string         `gorm:"size:30" json:"amount"`
	Fee         string         `gorm:"size:30" json:"fee"`
	FeeCurrency string         `gorm:"size:20" json:"fee_currency"`
	Role        string         `gorm:"size:10" json:"role"` // taker, maker
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Trade) TableName() string {
	return "trades"
}

type StrategyLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	StrategyID uint      `gorm:"index" json:"strategy_id"`
	EventType  string    `gorm:"size:50" json:"event_type"` // order_placed, order_filled, error, etc.
	Message    string    `gorm:"size:500" json:"message"`
	Data       string    `gorm:"type:text" json:"data,omitempty"` // JSON 数据
	CreatedAt  time.Time `json:"created_at"`
}

func (StrategyLog) TableName() string {
	return "strategy_logs"
}

// ==================== Copy Trading Models ====================

// PlatformCookie 平台级 Gate.io Cookie 配置（管理员绑定）
type PlatformCookie struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Token        string         `gorm:"size:2048;not null" json:"-"` // 加密存储
	CsrfToken    string         `gorm:"size:256;not null" json:"-"`  // 加密存储
	Uid          string         `gorm:"size:50;not null" json:"uid"`
	ExpiresAt    *time.Time     `gorm:"index" json:"expires_at"`
	Status       string         `gorm:"size:20;default:'active'" json:"status"` // active/expired
	LastSyncedAt *time.Time     `json:"last_synced_at"`
	ErrorMsg     string         `gorm:"size:512" json:"error_msg"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (PlatformCookie) TableName() string {
	return "platform_cookies"
}

// CopyTrader 交易员信息
type CopyTrader struct {
	ID              uint        `gorm:"primarykey" json:"id"`
	TraderID        string      `gorm:"size:100;not null;uniqueIndex" json:"trader_id"`
	TraderName      string      `gorm:"size:100" json:"trader_name"`
	Avatar          string      `gorm:"size:512" json:"avatar"`
	Exchange        string      `gorm:"size:20;default:'gate'" json:"exchange"`
	Status          string      `gorm:"size:20;default:'running'" json:"status"`

	// 统计数据
	Cycle           string      `gorm:"size:20" json:"cycle"`                   // day/week/month/all
	TotalPnl        string      `gorm:"size:32" json:"total_pnl"`
	TotalRoi        string      `gorm:"size:32" json:"total_roi"`
	FollowProfit    string      `gorm:"size:32" json:"follow_profit"`
	FollowRoi       string      `gorm:"size:32" json:"follow_roi"`
	WinRate         string      `gorm:"size:16" json:"win_rate"`
	FollowerCount   int         `gorm:"default:0" json:"follower_count"`
	PositionCount   int         `gorm:"default:0" json:"position_count"`
	MaxDrawdown     string      `gorm:"size:32" json:"max_drawdown"`
	AvgLeverage     string      `gorm:"size:16" json:"avg_leverage"`

	// 属性
	IsCurated       bool        `gorm:"default:false" json:"is_curated"`
	IsPrivate       bool        `gorm:"default:false" json:"is_private"`
	StyleLabels     string      `gorm:"type:jsonb" json:"style_labels"` // JSON 数组

	// 同步
	LastSyncedAt    *time.Time  `gorm:"index" json:"last_synced_at"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

func (CopyTrader) TableName() string {
	return "copy_traders"
}

// CopyTraderDailyStats 交易员每日统计快照
type CopyTraderDailyStats struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	TraderID        string    `gorm:"size:100;not null;index" json:"trader_id"`
	Date            string    `gorm:"size:10;not null;index" json:"date"` // YYYY-MM-DD
	TotalPnl        string    `gorm:"size:32" json:"total_pnl"`
	TotalRoi        string    `gorm:"size:32" json:"total_roi"`
	FollowProfit    string    `gorm:"size:32" json:"follow_profit"`
	FollowerCount   int       `gorm:"default:0" json:"follower_count"`
	CreatedAt       time.Time `json:"created_at"`
}

func (CopyTraderDailyStats) TableName() string {
	return "copy_trader_daily_stats"
}

// AdminNotification 管理员通知
type AdminNotification struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Type      string         `gorm:"size:50;not null" json:"type"` // cookie_expiry, system, etc.
	Message   string         `gorm:"size:512;not null" json:"message"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AdminNotification) TableName() string {
	return "admin_notifications"
}
