
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/falser101/gate-trading/internal/models"
	"github.com/falser101/gate-trading/internal/repository"
	"github.com/falser101/gate-trading/internal/service/copytrading"
)

// ==================== Copy Trading Handlers ====================

// handleGetCopyTraders 获取交易员列表
func handleGetCopyTraders(svc *copytrading.SyncService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		orderBy := c.DefaultQuery("order_by", "follow_profit")
		sortBy := c.DefaultQuery("sort_by", "desc")
		cycle := c.DefaultQuery("cycle", "month")
		status := c.DefaultQuery("status", "running")
		searchKey := c.Query("search")

		result, err := svc.GetTraderList(page, pageSize, orderBy, sortBy, cycle, status, searchKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":      result.List,
			"total":     result.Total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

// handleGetCopyTraderDetail 获取交易员详情
func handleGetCopyTraderDetail(svc *copytrading.SyncService) gin.HandlerFunc {
	return func(c *gin.Context) {
		traderID := c.Param("id")

		trader, err := svc.GetTraderDetail(traderID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trader not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": trader})
	}
}

// handleGetCopyTraderStats 获取交易员历史统计
func handleGetCopyTraderStats(svc *copytrading.SyncService) gin.HandlerFunc {
	return func(c *gin.Context) {
		traderID := c.Param("id")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		stats, err := svc.GetTraderStats(traderID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": stats})
	}
}

// CopyCookieRequest Cookie 绑定请求
type CopyCookieRequest struct {
	Token     string `json:"token" binding:"required"`
	CsrfToken string `json:"csrftoken" binding:"required"`
	Uid       string `json:"uid" binding:"required"`
}

// handleSaveCookie 保存 Cookie
func handleSaveCookie(svc *copytrading.SyncService, cookieMgr *copytrading.CookieManager, repo *repository.CopyTradingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CopyCookieRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 解析 Cookie
		gateCookie, err := cookieMgr.NewGateCookie(req.Token, req.CsrfToken, req.Uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookie: " + err.Error()})
			return
		}

		// 加密
		encryptedToken, encryptedCsrf, err := cookieMgr.EncryptGateCookie(gateCookie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt cookie"})
			return
		}

		// 保存
		cookie := &models.PlatformCookie{
			Token:        encryptedToken,
			CsrfToken:    encryptedCsrf,
			Uid:          req.Uid,
			ExpiresAt:    &gateCookie.ExpiresAt,
			Status:       "active",
			LastSyncedAt: nil,
		}

		if err := repo.CreateOrUpdateCookie(cookie); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":           "Cookie saved successfully",
			"expires_at":        gateCookie.ExpiresAt.Format(time.RFC3339),
			"days_until_expiry": time.Until(gateCookie.ExpiresAt).Hours() / 24,
		})
	}
}

// handleGetCookie 获取 Cookie 状态
func handleGetCookie(repo *repository.CopyTradingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := repo.GetActiveCookie()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active cookie found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":             cookie.ID,
			"uid":            cookie.Uid,
			"status":         cookie.Status,
			"expires_at":     cookie.ExpiresAt,
			"last_synced_at": cookie.LastSyncedAt,
			"error_msg":      cookie.ErrorMsg,
			"created_at":     cookie.CreatedAt,
		})
	}
}

// handleDeleteCookie 删除 Cookie
func handleDeleteCookie(repo *repository.CopyTradingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := repo.GetActiveCookie()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active cookie found"})
			return
		}

		// 软删除
		if err := repo.DeleteCookie(cookie); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cookie deleted successfully"})
	}
}

// handleManualSync 手动触发同步
func handleManualSync(svc *copytrading.SyncService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := svc.SyncAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":         "Sync completed",
			"synced_count":    result.SyncedCount,
			"new_count":       result.NewCount,
			"updated_count":   result.UpdatedCount,
			"error_count":     result.ErrorCount,
			"duration_seconds": result.DurationSecond,
		})
	}
}

// handleGetNotifications 获取管理员通知
func handleGetNotifications(repo *repository.CopyTradingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		notifications, err := repo.GetAllNotifications(page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": notifications})
	}
}

// handleMarkNotificationRead 标记通知为已读
func handleMarkNotificationRead(repo *repository.CopyTradingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID uint `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.MarkNotificationAsRead(req.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
	}
}
