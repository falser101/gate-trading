package strategy

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	gateapi "github.com/gate/gateapi-go/v7"
	"github.com/falser101/gate-trading/internal/gateway"
	"github.com/falser101/gate-trading/internal/models"
	"github.com/falser101/gate-trading/internal/repository"
)

// 策略类型
const (
	StrategyTypeGrid = "grid"
	StrategyTypeDCA  = "dca"
)

// 策略状态
const (
	StatusRunning = "running"
	StatusStopped = "stopped"
	StatusPaused  = "paused"
)

// 策略接口
type Strategy interface {
	ID() uint
	Type() string
	OnTick(ctx context.Context, price decimal.Decimal) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Config() interface{}
}

// 策略引擎
type Engine struct {
	mu            sync.RWMutex
	strategies    map[uint]Strategy
	strategyRepos *repository.StrategyRepository
	orderRepos    *repository.OrderRepository
	logRepos      *repository.StrategyLogRepository
	gateFactory   *gateway.GateClientFactory
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewStrategyEngine(
	strategyRepo *repository.StrategyRepository,
	orderRepo *repository.OrderRepository,
	logRepo *repository.StrategyLogRepository,
	gateFactory *gateway.GateClientFactory,
) *Engine {
	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		strategies:    make(map[uint]Strategy),
		strategyRepos: strategyRepo,
		orderRepos:    orderRepo,
		logRepos:      logRepo,
		gateFactory:   gateFactory,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// 创建策略
func (e *Engine) CreateStrategy(userID uint, strategyType, symbol string, config map[string]interface{}) (*models.Strategy, error) {
	strategy := &models.Strategy{
		UserID: userID,
		Type:   strategyType,
		Symbol: symbol,
		Config: config,
		Status: StatusStopped,
	}

	// 保存策略
	if err := e.strategyRepos.Create(strategy); err != nil {
		return nil, err
	}

	return strategy, nil
}

// 启动策略
func (e *Engine) StartStrategy(ctx context.Context, id uint) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 检查策略是否已存在
	if _, exists := e.strategies[id]; exists {
		return fmt.Errorf("strategy %d is already running", id)
	}

	// 获取策略配置
	strategy, err := e.strategyRepos.GetByID(id)
	if err != nil {
		return err
	}

	// 获取用户信息（通过 strategy 中的 UserID）
	var userModel models.User
	if err := e.strategyRepos.GetDB().Model(&models.User{}).Where("id = ?", strategy.UserID).First(&userModel).Error; err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if userModel.GateAPIKey == "" || userModel.GateAPISecret == "" {
		return fmt.Errorf("user has not bound Gate API key")
	}

	client := e.gateFactory.GetClient(userModel.GateAPIKey, userModel.GateAPISecret)

	// 创建策略实例
	var s Strategy
	switch strategy.Type {
	case StrategyTypeGrid:
		s = NewGridStrategy(strategy, client, e)
	case StrategyTypeDCA:
		s = NewDCAStrategy(strategy, client, e)
	default:
		return fmt.Errorf("unknown strategy type: %s", strategy.Type)
	}

	// 启动策略
	if err := s.Start(ctx); err != nil {
		return err
	}

	// 更新策略状态
	if err := e.strategyRepos.UpdateStatus(id, StatusRunning); err != nil {
		return err
	}

	e.strategies[id] = s

	// 记录日志
	e.Log(id, "start", "Strategy started", nil)

	return nil
}

// 停止策略
func (e *Engine) StopStrategy(ctx context.Context, id uint) error {
	e.mu.Lock()
	strategy, exists := e.strategies[id]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("strategy %d is not running", id)
	}
	e.mu.Unlock()

	// 停止策略
	if err := strategy.Stop(ctx); err != nil {
		return err
	}

	// 更新策略状态
	if err := e.strategyRepos.UpdateStatus(id, StatusStopped); err != nil {
		return err
	}

	e.mu.Lock()
	delete(e.strategies, id)
	e.mu.Unlock()

	// 记录日志
	e.Log(id, "stop", "Strategy stopped", nil)

	return nil
}

// 删除策略
func (e *Engine) DeleteStrategy(id uint) error {
	// 如果策略正在运行，先停止
	ctx := context.Background()
	if s, exists := e.strategies[id]; exists {
		if err := s.Stop(ctx); err != nil {
			return err
		}
		delete(e.strategies, id)
	}

	return e.strategyRepos.Delete(id)
}

// 记录策略日志
func (e *Engine) Log(strategyID uint, eventType, message string, data map[string]interface{}) error {
	var dataStr string
	if data != nil {
		b, _ := json.Marshal(data)
		dataStr = string(b)
	}

	log := &models.StrategyLog{
		StrategyID: strategyID,
		EventType:  eventType,
		Message:    message,
		Data:       dataStr,
	}

	return e.logRepos.Create(log)
}

// 网格策略
type GridStrategy struct {
	strategy     *models.Strategy
	client       *gateapi.APIClient
	engine       *Engine
	grids        []GridLevel
	activeOrders map[string]uint // order_id -> grid_index
	mu           sync.RWMutex
}

type GridLevel struct {
	Index       int
	BuyPrice    decimal.Decimal
	SellPrice   decimal.Decimal
	Amount      decimal.Decimal
	BuyOrderID  string
	SellOrderID string
	Filled      bool
}

func NewGridStrategy(s *models.Strategy, client *gateapi.APIClient, engine *Engine) *GridStrategy {
	return &GridStrategy{
		strategy:     s,
		client:       client,
		engine:       engine,
		activeOrders: make(map[string]uint),
	}
}

func (g *GridStrategy) ID() uint                    { return g.strategy.ID }
func (g *GridStrategy) Type() string                { return StrategyTypeGrid }
func (g *GridStrategy) Config() interface{}         { return g.strategy.Config }
func (g *GridStrategy) OnTick(ctx context.Context, price decimal.Decimal) error { return nil }

func (g *GridStrategy) Start(ctx context.Context) error {
	// TODO: 实现网格策略启动逻辑
	// 需要实现：
	// 1. 解析配置
	// 2. 初始化网格
	// 3. 使用 g.client.SpotApi.GetTicker 获取当前价格
	// 4. 使用 g.client.SpotApi.CreateOrder 放置订单
	return nil
}

func (g *GridStrategy) Stop(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// TODO: 实现撤销订单逻辑
	// 需要使用 g.client.SpotApi.CancelOrder

	g.activeOrders = make(map[string]uint)
	return nil
}

// DCA 策略
type DCAStrategy struct {
	strategy *models.Strategy
	client   *gateapi.APIClient
	engine   *Engine
	timer    *time.Timer
	mu       sync.Mutex
}

func NewDCAStrategy(s *models.Strategy, client *gateapi.APIClient, engine *Engine) *DCAStrategy {
	return &DCAStrategy{
		strategy: s,
		client:   client,
		engine:   engine,
	}
}

func (d *DCAStrategy) ID() uint                    { return d.strategy.ID }
func (d *DCAStrategy) Type() string                { return StrategyTypeDCA }
func (d *DCAStrategy) Config() interface{}         { return d.strategy.Config }
func (d *DCAStrategy) OnTick(ctx context.Context, price decimal.Decimal) error { return nil }

func (d *DCAStrategy) Start(ctx context.Context) error {
	// TODO: 实现 DCA 策略启动逻辑
	return nil
}

func (d *DCAStrategy) Stop(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	return nil
}
