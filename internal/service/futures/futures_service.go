package futures

import (
	"context"

	gateapi "github.com/gate/gateapi-go/v7"
)

// FuturesService 合约服务
type FuturesService struct {
	client *gateapi.APIClient
	settle string // 结算币种，如 "usdt"
}

// NewFuturesService 创建合约服务
func NewFuturesService(client *gateapi.APIClient, settle string) *FuturesService {
	return &FuturesService{
		client: client,
		settle: settle,
	}
}

// GetAccount 获取合约账户信息
func (s *FuturesService) GetAccount(ctx context.Context) (*gateapi.FuturesAccount, error) {
	account, _, err := s.client.FuturesApi.ListFuturesAccounts(ctx, s.settle)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetPositions 获取持仓列表
func (s *FuturesService) GetPositions(ctx context.Context) ([]gateapi.Position, error) {
	positions, _, err := s.client.FuturesApi.ListPositions(ctx, s.settle, nil)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

// GetPosition 获取单个持仓
func (s *FuturesService) GetPosition(ctx context.Context, contract string) (*gateapi.Position, error) {
	position, _, err := s.client.FuturesApi.GetPosition(ctx, s.settle, contract)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// CreateOrder 创建订单
func (s *FuturesService) CreateOrder(ctx context.Context, order gateapi.FuturesOrder) (*gateapi.FuturesOrder, error) {
	createdOrder, _, err := s.client.FuturesApi.CreateFuturesOrder(ctx, s.settle, order, nil)
	if err != nil {
		return nil, err
	}
	return &createdOrder, nil
}

// GetOrders 获取订单列表
func (s *FuturesService) GetOrders(ctx context.Context, status string, limit int32) ([]gateapi.FuturesOrder, error) {
	opts := &gateapi.ListFuturesOrdersOpts{
		Limit: gateapi.NewOptionalInt32(limit),
	}
	if status != "" {
		opts.Status = gateapi.NewOptionalString(status)
	}
	orders, _, err := s.client.FuturesApi.ListFuturesOrders(ctx, s.settle, status, opts)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// CancelOrder 取消订单
func (s *FuturesService) CancelOrder(ctx context.Context, orderId string) (*gateapi.FuturesOrder, error) {
	cancelledOrder, _, err := s.client.FuturesApi.CancelFuturesOrder(ctx, s.settle, orderId, nil)
	if err != nil {
		return nil, err
	}
	return &cancelledOrder, nil
}

// CancelAllOrders 取消所有订单
func (s *FuturesService) CancelAllOrders(ctx context.Context) ([]gateapi.FuturesOrder, error) {
	orders, _, err := s.client.FuturesApi.CancelFuturesOrders(ctx, s.settle, nil)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// SetLeverage 调整杠杆
func (s *FuturesService) SetLeverage(ctx context.Context, contract string, leverage string, marginMode string) (*gateapi.Position, error) {
	position, _, err := s.client.FuturesApi.UpdateContractPositionLeverage(
		ctx, s.settle, contract, leverage, marginMode, nil)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// ClosePosition 平仓
func (s *FuturesService) ClosePosition(ctx context.Context, contract string) (*gateapi.FuturesOrder, error) {
	// 创建平仓订单：close=true, size=0
	order := gateapi.FuturesOrder{
		Contract: contract,
		Size:     "0",
		Price:    "0",
		Close:    true,
		Tif:      "ioc",
	}
	createdOrder, _, err := s.client.FuturesApi.CreateFuturesOrder(ctx, s.settle, order, nil)
	if err != nil {
		return nil, err
	}
	return &createdOrder, nil
}

// GetOrder 获取订单详情
func (s *FuturesService) GetOrder(ctx context.Context, orderId string) (*gateapi.FuturesOrder, error) {
	order, _, err := s.client.FuturesApi.GetFuturesOrder(ctx, s.settle, orderId)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetTickers 获取行情
func (s *FuturesService) GetTickers(ctx context.Context) ([]gateapi.FuturesTicker, error) {
	tickers, _, err := s.client.FuturesApi.ListFuturesTickers(ctx, s.settle, nil)
	if err != nil {
		return nil, err
	}
	return tickers, nil
}

// GetTicker 获取单个行情
func (s *FuturesService) GetTicker(ctx context.Context, contract string) (*gateapi.FuturesTicker, error) {
	ticker, _, err := s.client.FuturesApi.GetFuturesTicker(ctx, s.settle, contract)
	if err != nil {
		return nil, err
	}
	return &ticker, nil
}

// GetContracts 获取合约列表
func (s *FuturesService) GetContracts(ctx context.Context) ([]gateapi.Contract, error) {
	contracts, _, err := s.client.FuturesApi.ListFuturesContracts(ctx, s.settle, nil)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

// GetContract 获取合约详情
func (s *FuturesService) GetContract(ctx context.Context, contract string) (*gateapi.Contract, error) {
	contractInfo, _, err := s.client.FuturesApi.GetFuturesContract(ctx, s.settle, contract)
	if err != nil {
		return nil, err
	}
	return &contractInfo, nil
}
