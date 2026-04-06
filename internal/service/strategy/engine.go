package strategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/falser101/gate-trading/internal/gateway"
	"github.com/falser101/gate-trading/internal/models"
	"github.com/falser101/gate-trading/internal/repository"
	"github.com/falser101/gate-trading/pkg/gateapi"
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
	strategyRepos *repository.StrategyRepository,
	orderRepos *repository.OrderRepository,
	logRepos *repository.StrategyLogRepository,
	gateFactory *gateway.GateClientFactory,
) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		strategies:    make(map[uint]Strategy),
		strategyRepos: strategyRepos,
		orderRepos:    orderRepos,
		logRepos:      logRepos,
		gateFactory:  gateFactory,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// 加载用户的所有策略
func (e *Engine) LoadUserStrategies(userID uint, user *models.User) error {
	strategies, err := e.strategyRepos.GetByUserID(userID)
	if err != nil {
		return err
	}

	for _, s := range strategies {
		if s.Status == StatusRunning {
			if err := e.LoadStrategy(&s, user); err != nil {
				return err
			}
		}
	}
	return nil
}

// 加载单个策略
func (e *Engine) LoadStrategy(s *models.Strategy, user *models.User) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	client := e.gateFactory.GetClient(user.GateAPIKey, user.GateAPISecret)

	var strategy Strategy
	switch s.Type {
	case StrategyTypeGrid:
		strategy = NewGridStrategy(s, client, e)
	case StrategyTypeDCA:
		strategy = NewDCAStrategy(s, client, e)
	default:
		return errors.New("unknown strategy type")
	}

	e.strategies[s.ID] = strategy
	return nil
}

// 启动策略
func (e *Engine) StartStrategy(id uint) error {
	e.mu.RLock()
	strategy, ok := e.strategies[id]
	e.mu.RUnlock()

	if !ok {
		return errors.New("strategy not found")
	}

	if err := e.strategyRepos.UpdateStatus(id, StatusRunning); err != nil {
		return err
	}

	return strategy.Start(e.ctx)
}

// 停止策略
func (e *Engine) StopStrategy(id uint) error {
	e.mu.RLock()
	strategy, ok := e.strategies[id]
	e.mu.RUnlock()

	if !ok {
		return errors.New("strategy not found")
	}

	if err := e.strategyRepos.UpdateStatus(id, StatusStopped); err != nil {
		return err
	}

	return strategy.Stop(e.ctx)
}

// 创建策略
func (e *Engine) CreateStrategy(userID uint, sType, symbol string, cfg interface{}) (*models.Strategy, error) {
	configMap, ok := cfg.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid config type")
	}

	strategy := &models.Strategy{
		UserID:   userID,
		Type:     sType,
		Symbol:   symbol,
		Config:   configMap,
		Status:   StatusStopped,
		RunState: models.RunStateJSON{},
	}

	if err := e.strategyRepos.Create(strategy); err != nil {
		return nil, err
	}

	return strategy, nil
}

// GetStrategy 获取策略实例
func (e *Engine) GetStrategy(id uint) Strategy {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.strategies[id]
}

// 记录日志
func (e *Engine) Log(strategyID uint, eventType, message string, data interface{}) error {
	dataStr := ""
	if data != nil {
		if b, err := json.Marshal(data); err == nil {
			dataStr = string(b)
		}
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
	client       *gateapi.Client
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

func NewGridStrategy(s *models.Strategy, client *gateapi.Client, engine *Engine) *GridStrategy {
	return &GridStrategy{
		strategy:     s,
		client:       client,
		engine:       engine,
		activeOrders: make(map[string]uint),
	}
}

func (g *GridStrategy) ID() uint                     { return g.strategy.ID }
func (g *GridStrategy) Type() string                 { return StrategyTypeGrid }
func (g *GridStrategy) Config() interface{}          { return g.strategy.Config }

func (g *GridStrategy) Start(ctx context.Context) error {
	// 解析配置
	cfg := g.strategy.Config
	var gridCfg models.GridConfig

	if v, ok := cfg["lower_price"].(string); ok {
		gridCfg.LowerPrice = v
	}
	if v, ok := cfg["upper_price"].(string); ok {
		gridCfg.UpperPrice = v
	}
	if v, ok := cfg["grid_count"].(float64); ok {
		gridCfg.GridCount = int(v)
	}
	if v, ok := cfg["invest_amount"].(string); ok {
		gridCfg.InvestAmount = v
	}
	if v, ok := cfg["profit_rate"].(string); ok {
		gridCfg.ProfitRate = v
	}

	// 初始化网格
	if err := g.initGrids(gridCfg); err != nil {
		return err
	}

	// 放置初始订单
	return g.placeInitialOrders(ctx)
}

func (g *GridStrategy) Stop(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 撤销所有活跃订单
	for orderID := range g.activeOrders {
		_, _ = g.client.CancelOrder(g.strategy.Symbol, orderID)
	}

	g.activeOrders = make(map[string]uint)
	return nil
}

func (g *GridStrategy) OnTick(ctx context.Context, price decimal.Decimal) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// 检查是否有订单成交
	for orderID, gridIndex := range g.activeOrders {
		order, err := g.client.GetOrder(g.strategy.Symbol, orderID)
		if err != nil {
			continue
		}

		if order.Status == "filled" {
			g.handleFilledOrder(ctx, gridIndex, order)
		}
	}

	return nil
}

func (g *GridStrategy) initGrids(cfg models.GridConfig) error {
	lowerPrice, _ := decimal.NewFromString(cfg.LowerPrice)
	upperPrice, _ := decimal.NewFromString(cfg.UpperPrice)
	investAmount, _ := decimal.NewFromString(cfg.InvestAmount)

	gridCount := decimal.NewFromInt(int64(cfg.GridCount))
	priceRange := upperPrice.Sub(lowerPrice)
	gridSize := priceRange.Div(gridCount)

	// 计算每格金额
	amountPerGrid := investAmount.Div(gridCount)

	g.grids = make([]GridLevel, cfg.GridCount)
	for i := 0; i < cfg.GridCount; i++ {
		buyPrice := lowerPrice.Add(gridSize.Mul(decimal.NewFromInt(int64(i))))
		sellPrice := buyPrice.Add(gridSize)
		amount := amountPerGrid.Div(buyPrice)

		g.grids[i] = GridLevel{
			Index:    i,
			BuyPrice: buyPrice,
			SellPrice: sellPrice,
			Amount:   amount,
			Filled:   false,
		}
	}

	return nil
}

func (g *GridStrategy) placeInitialOrders(ctx context.Context) error {
	// 获取当前价格
	ticker, err := g.client.GetTicker(g.strategy.Symbol)
	if err != nil {
		return err
	}

	currentPrice := ticker.Last

	// 在当前价格下方放置买单，上方放置卖单
	for i, grid := range g.grids {
		if grid.BuyPrice.LessThan(currentPrice) {
			// 放置买单
			resp, err := g.client.LimitBuy(g.strategy.Symbol, grid.Amount, grid.BuyPrice)
			if err != nil {
				g.engine.Log(g.strategy.ID, "order_error", fmt.Sprintf("Failed to place buy order: %v", err), nil)
				continue
			}

			g.mu.Lock()
			g.grids[i].BuyOrderID = resp.ID
			g.activeOrders[resp.ID] = uint(i)
			g.mu.Unlock()

			// 记录订单
			g.recordOrder(resp.ID, "buy", grid.BuyPrice, grid.Amount)
		}
	}

	return nil
}

func (g *GridStrategy) handleFilledOrder(ctx context.Context, gridIndex uint, order *gateapi.Order) {
	grid := g.grids[gridIndex]

	if order.Side == "buy" {
		// 买单成交，放置对应的卖单
		sellPrice := grid.BuyPrice.Mul(decimal.NewFromFloat(1.002)) // 0.2% 利润
		resp, err := g.client.LimitSell(g.strategy.Symbol, grid.Amount, sellPrice)
		if err != nil {
			g.engine.Log(g.strategy.ID, "order_error", fmt.Sprintf("Failed to place sell order: %v", err), nil)
			return
		}

		g.mu.Lock()
		g.grids[gridIndex].SellOrderID = resp.ID
		g.activeOrders[resp.ID] = gridIndex
		delete(g.activeOrders, order.ID)
		g.mu.Unlock()

		g.engine.Log(g.strategy.ID, "order_filled", fmt.Sprintf("Buy order filled at %s, sell order placed at %s",
			order.AveragePrice, sellPrice), nil)
	} else {
		// 卖单成交，重新放置买单
		resp, err := g.client.LimitBuy(g.strategy.Symbol, grid.Amount, grid.BuyPrice)
		if err != nil {
			g.engine.Log(g.strategy.ID, "order_error", fmt.Sprintf("Failed to place buy order: %v", err), nil)
			return
		}

		g.mu.Lock()
		g.grids[gridIndex].BuyOrderID = resp.ID
		g.activeOrders[resp.ID] = gridIndex
		delete(g.activeOrders, order.ID)
		g.mu.Unlock()

		g.engine.Log(g.strategy.ID, "order_filled", fmt.Sprintf("Sell order filled, buy order placed at %s",
			grid.BuyPrice), nil)
	}
}

func (g *GridStrategy) recordOrder(orderID, side string, price, amount decimal.Decimal) {
	order := &models.Order{
		StrategyID: g.strategy.ID,
		UserID:     g.strategy.UserID,
		GateOrderID: orderID,
		Symbol:     g.strategy.Symbol,
		Side:       side,
		Type:       "limit",
		Price:      price.String(),
		Amount:     amount.String(),
		Status:     "open",
	}
	g.engine.orderRepos.Create(order)
}

// DCA 策略
type DCAStrategy struct {
	strategy   *models.Strategy
	client     *gateapi.Client
	engine     *Engine
	lastBuyTime *time.Time
	buyCount    int
	mu          sync.RWMutex
}

func NewDCAStrategy(s *models.Strategy, client *gateapi.Client, engine *Engine) *DCAStrategy {
	return &DCAStrategy{
		strategy: s,
		client:   client,
		engine:   engine,
	}
}

func (d *DCAStrategy) ID() uint            { return d.strategy.ID }
func (d *DCAStrategy) Type() string        { return StrategyTypeDCA }
func (d *DCAStrategy) Config() interface{} { return d.strategy.Config }

func (d *DCAStrategy) Start(ctx context.Context) error {
	// DCA 策略定时执行，这里只初始化
	cfg := d.strategy.Config
	var dcaCfg models.DCAConfig

	if v, ok := cfg["invest_amount"].(string); ok {
		dcaCfg.InvestAmount = v
	}
	if v, ok := cfg["interval"].(float64); ok {
		dcaCfg.Interval = int(v)
	}
	if v, ok := cfg["target_price"].(string); ok {
		dcaCfg.TargetPrice = v
	}
	if v, ok := cfg["max_buy_times"].(float64); ok {
		dcaCfg.MaxBuyTimes = int(v)
	}

	d.engine.Log(d.strategy.ID, "strategy_started", "DCA strategy started", nil)
	return nil
}

func (d *DCAStrategy) Stop(ctx context.Context) error {
	d.engine.Log(d.strategy.ID, "strategy_stopped", "DCA strategy stopped", nil)
	return nil
}

func (d *DCAStrategy) OnTick(ctx context.Context, price decimal.Decimal) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	cfg := d.strategy.Config
	var dcaCfg models.DCAConfig

	if v, ok := cfg["invest_amount"].(string); ok {
		dcaCfg.InvestAmount = v
	}
	if v, ok := cfg["interval"].(float64); ok {
		dcaCfg.Interval = int(v)
	}
	if v, ok := cfg["target_price"].(string); ok {
		dcaCfg.TargetPrice = v
	}
	if v, ok := cfg["max_buy_times"].(float64); ok {
		dcaCfg.MaxBuyTimes = int(v)
	}

	// 检查是否达到买入间隔
	if d.lastBuyTime != nil && time.Since(*d.lastBuyTime) < time.Duration(dcaCfg.Interval)*time.Minute {
		return nil
	}

	// 检查是否达到最大买入次数
	if d.buyCount >= dcaCfg.MaxBuyTimes {
		return nil
	}

	// 检查目标价格（如果设置了）
	if dcaCfg.TargetPrice != "" {
		targetPrice, _ := decimal.NewFromString(dcaCfg.TargetPrice)
		if price.GreaterThan(targetPrice) {
			return nil // 当前价格高于目标价，不买入
		}
	}

	// 执行买入
	investAmount, _ := decimal.NewFromString(dcaCfg.InvestAmount)
	resp, err := d.client.MarketBuy(d.strategy.Symbol, investAmount)
	if err != nil {
		d.engine.Log(d.strategy.ID, "order_error", fmt.Sprintf("DCA buy failed: %v", err), nil)
		return err
	}

	d.buyCount++
	now := time.Now()
	d.lastBuyTime = &now

	d.engine.Log(d.strategy.ID, "order_placed", fmt.Sprintf("DCA buy #%d executed at %s", d.buyCount, price),
		map[string]string{"order_id": resp.ID})

	// 记录订单
	d.recordOrder(resp.ID, investAmount)

	return nil
}

func (d *DCAStrategy) recordOrder(orderID string, amount decimal.Decimal) {
	order := &models.Order{
		StrategyID:  d.strategy.ID,
		UserID:      d.strategy.UserID,
		GateOrderID: orderID,
		Symbol:      d.strategy.Symbol,
		Side:        "buy",
		Type:        "market",
		Amount:      amount.String(),
		Status:      "open",
	}
	d.engine.orderRepos.Create(order)
}
