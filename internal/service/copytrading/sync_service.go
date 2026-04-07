package copytrading

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/falser101/gate-trading/internal/models"
	"github.com/falser101/gate-trading/internal/repository"
)

// SyncService 跟单交易同步服务
type SyncService struct {
	repo      *repository.CopyTradingRepository
	cookieMgr *CookieManager
	client    *GateWebClient
}

// SyncResult 同步结果
type SyncResult struct {
	SyncedCount    int       `json:"synced_count"`
	UpdatedCount   int       `json:"updated_count"`
	NewCount       int       `json:"new_count"`
	ErrorCount     int       `json:"error_count"`
	SyncedAt       time.Time `json:"synced_at"`
	DurationSecond float64   `json:"duration_second"`
	Error          error     `json:"error"`
}

// NewSyncService 创建同步服务
func NewSyncService(
	repo *repository.CopyTradingRepository,
	cookieMgr *CookieManager,
	client *GateWebClient,
) *SyncService {
	return &SyncService{
		repo:      repo,
		cookieMgr: cookieMgr,
		client:    client,
	}
}

// SyncAll 同步所有交易员数据
func (s *SyncService) SyncAll() (*SyncResult, error) {
	startTime := time.Now()
	result := &SyncResult{SyncedAt: startTime}

	// 获取激活的 Cookie
	cookieRecord, err := s.repo.GetActiveCookie()
	if err != nil {
		result.Error = fmt.Errorf("failed to get active cookie: %w", err)
		return result, result.Error
	}

	// 解密 Cookie
	gateCookie, err := s.cookieMgr.DecryptGateCookie(
		cookieRecord.Token,
		cookieRecord.CsrfToken,
		cookieRecord.Uid,
		cookieRecord.ExpiresAt,
	)
	if err != nil {
		result.Error = fmt.Errorf("failed to decrypt cookie: %w", err)
		s.repo.UpdateCookieStatus(cookieRecord.ID, "error", result.Error.Error())
		return result, result.Error
	}

	// 检查 Cookie 是否过期
	if s.cookieMgr.IsExpired(gateCookie) {
		result.Error = fmt.Errorf("cookie has expired at %s", gateCookie.ExpiresAt.Format(time.RFC3339))
		s.repo.UpdateCookieStatus(cookieRecord.ID, "expired", result.Error.Error())
		return result, result.Error
	}

	// 检查是否即将过期（7 天内）
	if s.cookieMgr.IsExpiringSoon(gateCookie, 7) {
		// 创建管理员通知
		s.repo.CreateAdminNotification(&models.AdminNotification{
			Type:    "cookie_expiry",
			Message: fmt.Sprintf("Cookie 即将于 %s 过期，请尽快重新绑定", gateCookie.ExpiresAt.Format("2006-01-02 15:04")),
		})
	}

	// 获取所有交易员数据
	params := DefaultTraderListParams()
	allTraders, err := s.client.GetAllTraders(gateCookie, params)
	if err != nil {
		result.Error = fmt.Errorf("failed to fetch traders from API: %w", err)
		s.repo.UpdateCookieStatus(cookieRecord.ID, "error", result.Error.Error())
		return result, result.Error
	}

	// 处理交易员数据
	processResult, err := s.processTraders(cookieRecord.ID, allTraders)
	if err != nil {
		result.Error = err
		return result, err
	}

	// 更新同步结果
	result.SyncedCount = processResult.SyncedCount
	result.NewCount = processResult.NewCount
	result.UpdatedCount = processResult.UpdatedCount
	result.DurationSecond = time.Since(startTime).Seconds()

	// 更新 Cookie 最后同步时间
	s.repo.UpdateCookieSyncTime(cookieRecord.ID)

	return result, nil
}

// processTraders 处理交易员数据
func (s *SyncService) processTraders(cookieID uint, traders []TraderInfo) (*SyncResult, error) {
	result := &SyncResult{SyncedAt: time.Now()}

	now := time.Now()
	today := now.Format("2006-01-02")

	for _, t := range traders {
		// 获取交易员名称（优先使用 nickname，其次使用 hide_name）
		traderName := t.UserInfo.Nickname
		if traderName == "" {
			traderName = t.UserInfo.HideName
		}

		// 处理风格标签
		var styleLabels string
		if len(t.LabelInfo.Text) > 0 {
			labels := make([]string, 0, len(t.LabelInfo.Text))
			for _, label := range t.LabelInfo.Text {
				labels = append(labels, label.LabelName)
			}
			// 格式化为 JSON 数组字符串
			styleLabelsJSON, _ := json.Marshal(labels)
			styleLabels = string(styleLabelsJSON)
		}

		trader := &models.CopyTrader{
			TraderID:        fmt.Sprintf("%d", t.LeaderID),
			TraderName:      traderName,
			Avatar:          t.UserInfo.Avatar,
			Exchange:        "gate",
			Status:          "running",
			Cycle:           "month",
			TotalPnl:        fmt.Sprintf("%.2f", float64(t.Profit)),
			TotalRoi:        fmt.Sprintf("%.2f", float64(t.ProfitRate)*100), // profit_rate 是小数，转为百分比
			FollowProfit:    fmt.Sprintf("%.2f", float64(t.FollowProfit)),
			FollowRoi:       "",
			WinRate:         fmt.Sprintf("%.2f", float64(t.WinRate)*100), // win_rate 是小数，转为百分比
			FollowerCount:   t.CurrFollowNum,
			PositionCount:   0,
			MaxDrawdown:     fmt.Sprintf("%.2f", float64(t.MaxDrawdown)*100), // max_drawdown 是小数，转为百分比
			AvgLeverage:     "",
			IsCurated:       t.IsCurated,
			IsPrivate:       t.IsPrivateLeader,
			StyleLabels:     styleLabels,
			LastSyncedAt:    &now,
		}

		// 创建或更新交易员
		if err := s.repo.UpsertTrader(trader); err != nil {
			result.ErrorCount++
			continue
		}

		result.SyncedCount++

		// 创建每日统计快照
		dailyStats := &models.CopyTraderDailyStats{
			TraderID:      fmt.Sprintf("%d", t.LeaderID),
			Date:          today,
			TotalPnl:      fmt.Sprintf("%.2f", float64(t.Profit)),
			TotalRoi:      fmt.Sprintf("%.2f", float64(t.ProfitRate)*100),
			FollowProfit:  fmt.Sprintf("%.2f", float64(t.FollowProfit)),
			FollowerCount: t.CurrFollowNum,
		}
		if err := s.repo.CreateOrUpdateDailyStats(dailyStats); err != nil {
			// 统计失败不影响主流程
		}
	}

	return result, nil
}

// CheckCookieExpiry 检查 Cookie 过期状态
// 返回：(是否需要通知，过期时间，错误)
func (s *SyncService) CheckCookieExpiry() (bool, *time.Time, error) {
	cookieRecord, err := s.repo.GetActiveCookie()
	if err != nil {
		return false, nil, err
	}

	gateCookie, err := s.cookieMgr.DecryptGateCookie(
		cookieRecord.Token,
		cookieRecord.CsrfToken,
		cookieRecord.Uid,
		cookieRecord.ExpiresAt,
	)
	if err != nil {
		return false, nil, err
	}

	// 已过期
	if s.cookieMgr.IsExpired(gateCookie) {
		return true, &gateCookie.ExpiresAt, nil
	}

	// 即将过期（7 天内）
	if s.cookieMgr.IsExpiringSoon(gateCookie, 7) {
		return true, &gateCookie.ExpiresAt, nil
	}

	return false, &gateCookie.ExpiresAt, nil
}

// TestCookie 测试 Cookie 是否有效
func (s *SyncService) TestCookie(cookieID uint) (bool, int, error) {
	cookieRecord, err := s.repo.GetCookieByID(cookieID)
	if err != nil {
		return false, 0, err
	}

	gateCookie, err := s.cookieMgr.DecryptGateCookie(
		cookieRecord.Token,
		cookieRecord.CsrfToken,
		cookieRecord.Uid,
		cookieRecord.ExpiresAt,
	)
	if err != nil {
		return false, 0, err
	}

	// 尝试获取交易员列表
	params := DefaultTraderListParams()
	params.Page = 1
	params.PageSize = 1

	result, err := s.client.GetTraderList(gateCookie, params)
	if err != nil {
		return false, 0, err
	}

	return true, result.Data.TotalCount, nil
}

// FormatSyncResult 格式化同步结果为可读字符串
func (r *SyncResult) String() string {
	if r.Error != nil {
		return fmt.Sprintf("Sync failed: %v", r.Error)
	}
	return fmt.Sprintf(
		"Sync completed: %d total, %d new, %d updated, %d errors (%.2fs)",
		r.SyncedCount, r.NewCount, r.UpdatedCount, r.ErrorCount, r.DurationSecond,
	)
}

// GetTraderList 获取交易员列表（供 API 使用）
func (s *SyncService) GetTraderList(page, pageSize int, orderBy, sortBy, cycle, status, searchKey string) (*repository.TraderListResult, error) {
	params := &repository.TraderListParams{
		Page:      page,
		PageSize:  pageSize,
		OrderBy:   orderBy,
		SortBy:    sortBy,
		Cycle:     cycle,
		Status:    status,
		SearchKey: searchKey,
	}
	return s.repo.GetTraderList(params)
}

// GetTraderDetail 获取交易员详情
func (s *SyncService) GetTraderDetail(traderID string) (*models.CopyTrader, error) {
	return s.repo.GetTraderByID(traderID)
}

// GetTraderStats 获取交易员历史统计
func (s *SyncService) GetTraderStats(traderID, startDate, endDate string) ([]models.CopyTraderDailyStats, error) {
	return s.repo.GetTraderStats(traderID, startDate, endDate)
}
