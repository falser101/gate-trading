package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/falser101/gate-trading/internal/models"
)

type CopyTradingRepository struct {
	db *gorm.DB
}

func NewCopyTradingRepository(db *gorm.DB) *CopyTradingRepository {
	return &CopyTradingRepository{db: db}
}

// handleDeleteCookie 删除 Cookie 中使用了未导出的 db 字段，需要添加一个公开方法
// DeleteCookie 删除 Cookie
func (r *CopyTradingRepository) DeleteCookie(cookie *models.PlatformCookie) error {
	return r.db.Delete(cookie).Error
}

// ==================== PlatformCookie ====================

// GetActiveCookie 获取当前激活的 Cookie
func (r *CopyTradingRepository) GetActiveCookie() (*models.PlatformCookie, error) {
	var cookie models.PlatformCookie
	err := r.db.Where("status = ?", "active").First(&cookie).Error
	if err != nil {
		return nil, err
	}
	return &cookie, nil
}

// GetCookieByID 根据 ID 获取 Cookie
func (r *CopyTradingRepository) GetCookieByID(id uint) (*models.PlatformCookie, error) {
	var cookie models.PlatformCookie
	err := r.db.First(&cookie, id).Error
	if err != nil {
		return nil, err
	}
	return &cookie, nil
}

// CreateOrUpdateCookie 创建或更新 Cookie（单例模式）
func (r *CopyTradingRepository) CreateOrUpdateCookie(cookie *models.PlatformCookie) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先检查是否已存在
		var existing models.PlatformCookie
		err := tx.Where("status = ?", "active").First(&existing).Error
		if err == nil {
			// 已存在，更新
			existing.Token = cookie.Token
			existing.CsrfToken = cookie.CsrfToken
			existing.Uid = cookie.Uid
			existing.ExpiresAt = cookie.ExpiresAt
			existing.Status = cookie.Status
			return tx.Save(&existing).Error
		}

		// 不存在，创建
		return tx.Create(cookie).Error
	})
}

// UpdateCookieStatus 更新 Cookie 状态
func (r *CopyTradingRepository) UpdateCookieStatus(id uint, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return r.db.Model(&models.PlatformCookie{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateCookieSyncTime 更新 Cookie 最后同步时间
func (r *CopyTradingRepository) UpdateCookieSyncTime(id uint) error {
	now := time.Now()
	return r.db.Model(&models.PlatformCookie{}).Where("id = ?", id).Update("last_synced_at", now).Error
}

// ==================== CopyTrader ====================

// TraderListParams 交易员列表查询参数
type TraderListParams struct {
	Page      int
	PageSize  int
	OrderBy   string
	SortBy    string
	Cycle     string
	Status    string
	SearchKey string // 交易员名称搜索
}

// TraderListResult 交易员列表结果
type TraderListResult struct {
	List  []models.CopyTrader
	Total int64
}

// GetTraderList 获取交易员列表（分页）
func (r *CopyTradingRepository) GetTraderList(params *TraderListParams) (*TraderListResult, error) {
	var list []models.CopyTrader
	var total int64

	query := r.db.Model(&models.CopyTrader{})

	// 过滤条件
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Cycle != "" {
		query = query.Where("cycle = ?", params.Cycle)
	}
	if params.SearchKey != "" {
		query = query.Where("trader_name LIKE ?", "%"+params.SearchKey+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	orderBy := "follow_profit"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}
	sortBy := "DESC"
	if params.SortBy == "asc" {
		sortBy = "ASC"
	}
	query = query.Order(orderBy + " " + sortBy)

	// 分页
	offset := 0
	if params.Page > 1 {
		offset = (params.Page - 1) * params.PageSize
	}
	if err := query.Offset(offset).Limit(params.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	return &TraderListResult{List: list, Total: total}, nil
}

// GetTraderByID 根据交易员 ID 获取详情
func (r *CopyTradingRepository) GetTraderByID(traderID string) (*models.CopyTrader, error) {
	var trader models.CopyTrader
	err := r.db.Where("trader_id = ?", traderID).First(&trader).Error
	if err != nil {
		return nil, err
	}
	return &trader, nil
}

// UpsertTrader 创建或更新交易员信息
func (r *CopyTradingRepository) UpsertTrader(trader *models.CopyTrader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.CopyTrader
		err := tx.Where("trader_id = ?", trader.TraderID).First(&existing).Error
		if err == nil {
			// 已存在，更新
			trader.ID = existing.ID
			return tx.Save(trader).Error
		}

		// 不存在，创建
		return tx.Create(trader).Error
	})
}

// UpsertTraders 批量创建或更新交易员
func (r *CopyTradingRepository) UpsertTraders(traders []*models.CopyTrader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, t := range traders {
			if err := r.UpsertTrader(t); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteStaleTraders 删除长时间未更新的交易员（可选）
func (r *CopyTradingRepository) DeleteStaleTraders(threshold time.Time) error {
	return r.db.Where("last_synced_at < ? AND updated_at < ?", threshold, threshold).Delete(&models.CopyTrader{}).Error
}

// ==================== CopyTraderDailyStats ====================

// GetTraderStats 获取交易员历史统计
func (r *CopyTradingRepository) GetTraderStats(traderID string, startDate, endDate string) ([]models.CopyTraderDailyStats, error) {
	var stats []models.CopyTraderDailyStats
	query := r.db.Where("trader_id = ?", traderID)
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	query = query.Order("date ASC")

	err := query.Find(&stats).Error
	return stats, err
}

// GetTraderStatsByDate 根据日期获取统计
func (r *CopyTradingRepository) GetTraderStatsByDate(traderID, date string) (*models.CopyTraderDailyStats, error) {
	var stats models.CopyTraderDailyStats
	err := r.db.Where("trader_id = ? AND date = ?", traderID, date).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// CreateOrUpdateDailyStats 创建或更新每日统计
func (r *CopyTradingRepository) CreateOrUpdateDailyStats(stats *models.CopyTraderDailyStats) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.CopyTraderDailyStats
		err := tx.Where("trader_id = ? AND date = ?", stats.TraderID, stats.Date).First(&existing).Error
		if err == nil {
			stats.ID = existing.ID
			return tx.Save(stats).Error
		}
		return tx.Create(stats).Error
	})
}

// ==================== AdminNotification ====================

// CreateAdminNotification 创建管理员通知
func (r *CopyTradingRepository) CreateAdminNotification(notification *models.AdminNotification) error {
	return r.db.Create(notification).Error
}

// GetUnreadNotifications 获取未读通知
func (r *CopyTradingRepository) GetUnreadNotifications() ([]models.AdminNotification, error) {
	var notifications []models.AdminNotification
	err := r.db.Where("is_read = ?", false).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

// GetAllNotifications 获取所有通知（分页）
func (r *CopyTradingRepository) GetAllNotifications(page, pageSize int) ([]models.AdminNotification, error) {
	var notifications []models.AdminNotification
	offset := (page - 1) * pageSize
	err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error
	return notifications, err
}

// MarkNotificationAsRead 标记通知为已读
func (r *CopyTradingRepository) MarkNotificationAsRead(id uint) error {
	return r.db.Model(&models.AdminNotification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func (r *CopyTradingRepository) MarkAllNotificationsAsRead() error {
	return r.db.Model(&models.AdminNotification{}).Update("is_read", true).Error
}
