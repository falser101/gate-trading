package copytrading

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// FloatFromString 可以解析 float64 或 string 类型的数字
type FloatFromString float64

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (f *FloatFromString) UnmarshalJSON(data []byte) error {
	// 尝试直接解析为 float64
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FloatFromString(num)
		return nil
	}

	// 尝试解析为字符串
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*f = FloatFromString(v)
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FloatFromString", string(data))
}

// GateWebClient Gate.io Web API 客户端
type GateWebClient struct {
	httpClient *http.Client
	userAgent  string
	baseURL    string
}

// NewGateWebClient 创建 Gate Web API 客户端
func NewGateWebClient(userAgent string) *GateWebClient {
	return &GateWebClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: userAgent,
		baseURL:   "https://www.gate.com",
	}
}

// TraderListParams 交易员列表查询参数
type TraderListParams struct {
	SubWebsiteID    string  `json:"sub_website_id"`
	OrderBy         string  `json:"order_by"`
	SortBy          string  `json:"sort_by"`
	IsCurated       int     `json:"is_curated"`
	PrivateType     int     `json:"private_type"`
	Cycle           string  `json:"cycle"`
	LabelIDs        string  `json:"label_ids"`
	Status          string  `json:"status"`
	StyleLabelCode  string  `json:"style_label_code"`
	MaxProfit       float64 `json:"max_profit"`
	MinProfit       float64 `json:"min_profit"`
	MaxFollowProfit float64 `json:"max_follow_profit"`
	MinFollowProfit float64 `json:"min_follow_profit"`
	Page            int     `json:"page"`
	PageSize        int     `json:"page_size"`
	IsFavorite      bool    `json:"is_favorite"`
}

// DefaultTraderListParams 默认查询参数
func DefaultTraderListParams() *TraderListParams {
	return &TraderListParams{
		SubWebsiteID:    "0",
		OrderBy:         "overall",
		SortBy:          "desc",
		IsCurated:       0,
		PrivateType:     0,
		Cycle:           "month",
		LabelIDs:        "",
		Status:          "running",
		StyleLabelCode:  "",
		MaxProfit:       5000000,
		MinProfit:       0,
		MaxFollowProfit: 5000000,
		MinFollowProfit: 0,
		Page:            1,
		PageSize:        100,
		IsFavorite:      false,
	}
}

// TraderListResponse API 响应结构
type TraderListResponse struct {
	Code      int        `json:"code"`
	Message   string     `json:"message"`
	Data      TraderList `json:"data"`
	Timestamp int64      `json:"timestamp"`
}

// TraderList 交易员列表数据
type TraderList struct {
	List       []TraderInfo `json:"list"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pagesize"`
	PageCount  int          `json:"pagecount"`
	TotalCount int          `json:"totalcount"`
}

// TraderInfo 交易员信息
type TraderInfo struct {
	LeaderID        int              `json:"leader_id"`
	Level           int              `json:"level"`
	Profit          FloatFromString  `json:"profit"`
	ProfitRate      FloatFromString  `json:"profit_rate"`
	WinRate         FloatFromString  `json:"win_rate"`
	MaxDrawdown     FloatFromString  `json:"max_drawdown"`
	FollowProfit    FloatFromString  `json:"follow_profit"`
	CurrFollowNum   int              `json:"curr_follow_num"`
	MaxFollowNum    int              `json:"max_follow_num"`
	AUM             FloatFromString  `json:"aum"`
	SharpRatio      FloatFromString  `json:"sharp_ratio"`
	IsFollow        bool             `json:"is_follow"`
	IsFull          bool             `json:"is_full"`
	IsCurated       bool             `json:"is_curated"`
	IsPrivateLeader bool             `json:"is_private_leader"`
	LeadingDays     int              `json:"leading_days"`
	UserInfo        TraderUserInfo   `json:"user_info"`
	LabelInfo       TraderLabelInfo  `json:"label_info"`
}

// TraderUserInfo 交易员用户信息
type TraderUserInfo struct {
	Nick      string `json:"nick"`
	Nickname  string `json:"nickname"`
	HideName  string `json:"hide_name"`
	Tier      int    `json:"tier"`
	Avatar    string `json:"avatar"`
	Anonymous string `json:"anonymous"`
}

// TraderLabelInfo 交易员标签信息
type TraderLabelInfo struct {
	Badge []interface{}  `json:"badge"`
	Text  []LabelText    `json:"text"`
}

// LabelText 标签文本
type LabelText struct {
	LabelID   int    `json:"label_id"`
	LabelName string `json:"label_name"`
	LabelDesc string `json:"label_desc"`
}

// YieldCurveResponse 收益曲线响应
type YieldCurveResponse struct {
	Code      int                 `json:"code"`
	Message   string              `json:"message"`
	Data      YieldCurveData      `json:"data"`
	Timestamp int64               `json:"timestamp"`
}

// YieldCurveData 收益曲线数据
type YieldCurveData struct {
	List []YieldCurveItem `json:"list"`
}

// YieldCurveItem 收益曲线项
type YieldCurveItem struct {
	LeaderID          int     `json:"leader_id"`
	UserID            int     `json:"user_id"`
	LeaderYieldCurves []YieldCurve `json:"leader_yield_curves"`
}

// YieldCurve 收益曲线
type YieldCurve struct {
	Profit        string  `json:"profit"`
	ProfitRate    string  `json:"profit_rate"`
	CurrentProfit string  `json:"current_profit"`
	TotalInvest   string  `json:"total_invest"`
	CreateTime    int64   `json:"create_time"`
}

// GetTraderList 获取交易员列表（使用 recommend_list 端点）
func (c *GateWebClient) GetTraderList(cookie *GateCookie, params *TraderListParams) (*TraderListResponse, error) {
	// 构建 URL 参数
	queryParams := url.Values{}
	queryParams.Set("sub_website_id", params.SubWebsiteID)
	queryParams.Set("order_by", params.OrderBy)
	queryParams.Set("sort_by", params.SortBy)
	queryParams.Set("is_curated", strconv.Itoa(params.IsCurated))
	queryParams.Set("private_type", strconv.Itoa(params.PrivateType))
	queryParams.Set("cycle", params.Cycle)
	queryParams.Set("label_ids", params.LabelIDs)
	queryParams.Set("status", params.Status)
	queryParams.Set("style_label_code", params.StyleLabelCode)

	if params.MaxProfit > 0 {
		queryParams.Set("max_profit", strconv.FormatFloat(params.MaxProfit, 'f', -1, 64))
	}
	if params.MinProfit > 0 {
		queryParams.Set("min_profit", strconv.FormatFloat(params.MinProfit, 'f', -1, 64))
	}
	if params.MaxFollowProfit > 0 {
		queryParams.Set("max_follow_profit", strconv.FormatFloat(params.MaxFollowProfit, 'f', -1, 64))
	}
	if params.MinFollowProfit > 0 {
		queryParams.Set("min_follow_profit", strconv.FormatFloat(params.MinFollowProfit, 'f', -1, 64))
	}

	queryParams.Set("page", strconv.Itoa(params.Page))
	queryParams.Set("page_size", strconv.Itoa(params.PageSize))

	urlStr := fmt.Sprintf("%s/apiw/v2/copy/leader/recommend_list?%s", c.baseURL, queryParams.Encode())

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Cookie", fmt.Sprintf("token=%s; csrftoken=%s; uid=%s", cookie.Token, cookie.CsrfToken, cookie.Uid))
	req.Header.Set("Csrftoken", cookie.CsrfToken)
	req.Header.Set("Sub-Website-Id", params.SubWebsiteID)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://www.gate.com/zh/copytrading")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TraderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: code=%d, message=%s", result.Code, result.Message)
	}

	return &result, nil
}

// GetYieldCurve 获取交易员收益曲线
func (c *GateWebClient) GetYieldCurve(cookie *GateCookie, leaderIDs []int) (*YieldCurveResponse, error) {
	urlStr := fmt.Sprintf("%s/apiw/v2/copy/api/leader/yield_curve?sub_website_id=0", c.baseURL)

	// 构建请求体
	body := map[string]interface{}{
		"leader_ids": leaderIDs,
		"data_type":  "month",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Cookie", fmt.Sprintf("token=%s; csrftoken=%s; uid=%s", cookie.Token, cookie.CsrfToken, cookie.Uid))
	req.Header.Set("Csrftoken", cookie.CsrfToken)
	req.Header.Set("Sub-Website-Id", "0")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.gate.com/zh/copytrading")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result YieldCurveResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: code=%d, message=%s", result.Code, result.Message)
	}

	return &result, nil
}

// GetAllTraders 获取所有交易员（自动分页）
func (c *GateWebClient) GetAllTraders(cookie *GateCookie, params *TraderListParams) ([]TraderInfo, error) {
	var allTraders []TraderInfo
	page := 1

	for {
		params.Page = page
		result, err := c.GetTraderList(cookie, params)
		if err != nil {
			return nil, err
		}

		allTraders = append(allTraders, result.Data.List...)

		// 如果已经获取完所有数据，退出循环
		if len(result.Data.List) < params.PageSize || page >= result.Data.PageCount {
			break
		}

		page++

		// 防止无限循环，最多获取 10 页
		if page > 10 {
			break
		}
	}

	return allTraders, nil
}
