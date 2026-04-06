package gateapi

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/shopspring/decimal"
)

// 获取账户余额
func (c *Client) GetSpotAccount() (map[string]Account, error) {
	var accounts []Account
	err := c.Get("/spot/accounts", &accounts)
	if err != nil {
		return nil, err
	}

	result := make(map[string]Account)
	for _, acc := range accounts {
		result[acc.Currency] = acc
	}
	return result, nil
}

// 获取单个币种余额
func (c *Client) GetBalance(currency string) (*Account, error) {
	accounts, err := c.GetSpotAccount()
	if err != nil {
		return nil, err
	}

	if acc, ok := accounts[currency]; ok {
		return &acc, nil
	}

	return &Account{
		Currency:  currency,
		Available: decimal.Zero,
		Locked:    decimal.Zero,
	}, nil
}

// 现货下单
func (c *Client) PlaceOrder(req *PlaceOrderRequest) (*PlaceOrderResponse, error) {
	var resp PlaceOrderResponse
	err := c.Post("/spot/orders", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 限价买单
func (c *Client) LimitBuy(currencyPair string, amount, price decimal.Decimal) (*PlaceOrderResponse, error) {
	return c.PlaceOrder(&PlaceOrderRequest{
		CurrencyPair: currencyPair,
		Type:         "limit",
		Side:         "buy",
		Amount:       amount.String(),
		Price:        price.String(),
		TimeInForce:  "gtc",
	})
}

// 限价卖单
func (c *Client) LimitSell(currencyPair string, amount, price decimal.Decimal) (*PlaceOrderResponse, error) {
	return c.PlaceOrder(&PlaceOrderRequest{
		CurrencyPair: currencyPair,
		Type:         "limit",
		Side:         "sell",
		Amount:       amount.String(),
		Price:        price.String(),
		TimeInForce:  "gtc",
	})
}

// 市价买单
func (c *Client) MarketBuy(currencyPair string, amount decimal.Decimal) (*PlaceOrderResponse, error) {
	return c.PlaceOrder(&PlaceOrderRequest{
		CurrencyPair: currencyPair,
		Type:         "market",
		Side:         "buy",
		Amount:       amount.String(),
	})
}

// 市价卖单
func (c *Client) MarketSell(currencyPair string, amount decimal.Decimal) (*PlaceOrderResponse, error) {
	return c.PlaceOrder(&PlaceOrderRequest{
		CurrencyPair: currencyPair,
		Type:         "market",
		Side:         "sell",
		Amount:       amount.String(),
	})
}

// 查询订单
func (c *Client) GetOrder(currencyPair, orderID string) (*Order, error) {
	path := fmt.Sprintf("/spot/orders/%s?currency_pair=%s", url.PathEscape(orderID), currencyPair)
	var order Order
	err := c.Get(path, &order)
	return &order, err
}

// 撤销订单
func (c *Client) CancelOrder(currencyPair, orderID string) (*Order, error) {
	path := fmt.Sprintf("/spot/orders/%s?currency_pair=%s", url.PathEscape(orderID), currencyPair)
	var order Order
	err := c.Delete(path, &order)
	return &order, err
}

// 查询当前委托订单
func (c *Client) GetOpenOrders(currencyPair string, offset, limit int) ([]Order, error) {
	query := url.Values{}
	query.Set("currency_pair", currencyPair)
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", strconv.Itoa(limit))

	path := "/spot/orders?" + query.Encode()
	var orders []Order
	err := c.Get(path, &orders)
	return orders, err
}

// 查询历史订单
func (c *Client) GetOrderHistory(currencyPair string, offset, limit int) ([]Order, error) {
	query := url.Values{}
	query.Set("currency_pair", currencyPair)
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", strconv.Itoa(limit))

	path := "/spot/order_history?" + query.Encode()
	var orders []Order
	err := c.Get(path, &orders)
	return orders, err
}

// 获取 Ticker
func (c *Client) GetTicker(currencyPair string) (*Ticker, error) {
	query := url.Values{}
	query.Set("currency_pair", currencyPair)

	path := "/spot/tickers?" + query.Encode()
	var tickers []Ticker
	err := c.Get(path, &tickers)
	if err != nil {
		return nil, err
	}

	if len(tickers) > 0 {
		return &tickers[0], nil
	}
	return nil, fmt.Errorf("no ticker found for %s", currencyPair)
}

// 获取 K 线数据
// interval: 1m, 5m, 15m, 30m, 1h, 4h, 1d, 7d, 30d
func (c *Client) GetCandlesticks(currencyPair, interval string, from, to int64, limit int) ([]Candlestick, error) {
	query := url.Values{}
	query.Set("currency_pair", currencyPair)
	query.Set("interval", interval)
	query.Set("from", strconv.FormatInt(from, 10))
	query.Set("to", strconv.FormatInt(to, 10))
	query.Set("limit", strconv.Itoa(limit))

	path := "/spot/candlesticks?" + query.Encode()
	var candles [][]interface{}
	err := c.Get(path, &candles)
	if err != nil {
		return nil, err
	}

	result := make([]Candlestick, 0, len(candles))
	for _, c := range candles {
		if len(c) >= 6 {
			candle := Candlestick{}
			candle.T, _ = c[0].(float64), nil
			candle.V, _ = decimal.NewFromString(fmt.Sprintf("%v", c[1]))
			candle.C, _ = decimal.NewFromString(fmt.Sprintf("%v", c[2]))
			candle.H, _ = decimal.NewFromString(fmt.Sprintf("%v", c[3]))
			candle.L, _ = decimal.NewFromString(fmt.Sprintf("%v", c[4]))
			candle.O, _ = decimal.NewFromString(fmt.Sprintf("%v", c[5]))
			result = append(result, candle)
		}
	}

	return result, nil
}

// 获取所有交易对
func (c *Client) GetCurrencyPairs() ([]map[string]interface{}, error) {
	var pairs []map[string]interface{}
	err := c.Get("/spot/currency_pairs", &pairs)
	return pairs, err
}
