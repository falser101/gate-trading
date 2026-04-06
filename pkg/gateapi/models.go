package gateapi

import "github.com/shopspring/decimal"

// 现货订单
type Order struct {
	ID             string          `json:"id"`
	Text           string          `json:"text"`
	CreateTime     int64           `json:"create_time"`
	CreateTimeMs   int64           `json:"create_time_ms"`
	UpdateTime     int64           `json:"update_time"`
	UpdateTimeMs   int64           `json:"update_time_ms"`
	Status         string          `json:"status"` // open, closed, cancelled
	CurrencyPair   string          `json:"currency_pair"`
	Type           string          `json:"type"` // limit, market
	Side           string          `json:"side"` // buy, sell
	Amount         decimal.Decimal `json:"amount"`
	Price          decimal.Decimal `json:"price"`
	TimeInForce    string          `json:"time_in_force"` // gtc, ioc
	AveragePrice   decimal.Decimal `json:"avg_deal_price"`
	FilledAmount   decimal.Decimal `json:"filled_amount"`
	FilledValue    decimal.Decimal `json:"filled_total"`
	RemainingAmount decimal.Decimal `json:"left_amount"`
	Fee            decimal.Decimal `json:"fee"`
	FeeCurrency    string          `json:"fee_currency"`
	FinishAs       string          `json:"finish_as"` // open, filled, cancelled
}

// 下单请求
type PlaceOrderRequest struct {
	CurrencyPair string          `json:"currency_pair"`
	Type         string          `json:"type"` // limit, market
	Side         string          `json:"side"` // buy, sell
	Amount       string          `json:"amount"`
	Price        string          `json:"price,omitempty"` // 市价单可不传
	TimeInForce  string          `json:"time_in_force,omitempty"` // gtc, ioc, fok
	AutoBorrow   bool            `json:"auto_borrow,omitempty"`
}

// 订单响应
type PlaceOrderResponse struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	CurrencyPair string `json:"currency_pair"`
}

// 账户余额
type Account struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
	UpdateID  int64           `json:"update_id"`
}

// Ticker
type Ticker struct {
	CurrencyPair string          `json:"currency_pair"`
	Last         decimal.Decimal `json:"last"`
	LowestAsk    decimal.Decimal `json:"lowest_ask"`
	HighestBid   decimal.Decimal `json:"highest_bid"`
	Change24h    decimal.Decimal `json:"change_percentage"`
	BaseVolume   decimal.Decimal `json:"base_volume"`
	QuoteVolume  decimal.Decimal `json:"quote_volume"`
	High24h      decimal.Decimal `json:"high_24h"`
	Low24h       decimal.Decimal `json:"low_24h"`
}

// K 线数据
type Candlestick struct {
	T         int64           `json:"t"` // 时间戳
	V         decimal.Decimal `json:"v"` // 成交量
	C         decimal.Decimal `json:"c"` // 收盘价
	H         decimal.Decimal `json:"h"` // 最高价
	L         decimal.Decimal `json:"l"` // 最低价
	O         decimal.Decimal `json:"o"` // 开盘价
	Sum       decimal.Decimal `json:"sum"`
}

// 交易记录
type Trade struct {
	ID           int64           `json:"id"`
	CreateTime   int64           `json:"create_time"`
	CurrencyPair string          `json:"currency_pair"`
	Side         string          `json:"side"`
	Amount       decimal.Decimal `json:"amount"`
	Price        decimal.Decimal `json:"price"`
	OrderID      string          `json:"order_id"`
	Role         string          `json:"role"` // taker, maker
	Fee          decimal.Decimal `json:"fee"`
	FeeCurrency  string          `json:"fee_currency"`
}
