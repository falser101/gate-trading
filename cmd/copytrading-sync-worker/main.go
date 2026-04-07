package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/falser101/gate-trading/internal/config"
	"github.com/falser101/gate-trading/internal/repository"
	"github.com/falser101/gate-trading/internal/models"
	"github.com/falser101/gate-trading/internal/service/copytrading"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 {
		handleCommand(os.Args[1:])
		return
	}

	// 默认运行定时同步任务
	runWorker()
}

// runWorker 运行定时同步任务
func runWorker() {
	log.Println("=== Copy Trading Sync Worker Starting ===")

	if err := config.Init(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}
	cfg := config.C

	db, err := repository.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Running database migration...")
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	repo := repository.NewCopyTradingRepository(db.DB)
	cookieMgr := copytrading.NewCookieManager(cfg.CopyTrading.EncryptionKey)
	webClient := copytrading.NewGateWebClient(cfg.CopyTrading.UserAgent)
	syncSvc := copytrading.NewSyncService(repo, cookieMgr, webClient)

	syncInterval := time.Duration(cfg.CopyTrading.SyncIntervalMinutes) * time.Minute
	if syncInterval < 5*time.Minute {
		syncInterval = 5 * time.Minute
	}
	log.Printf("Sync interval: %d minutes", int(syncInterval.Minutes()))

	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down...")
		cancel()
	}()

	// 首次同步
	log.Println("Running initial sync...")
	result, err := syncSvc.SyncAll()
	if err != nil {
		log.Printf("Initial sync failed: %v", err)
	} else {
		log.Printf("Initial sync: %d synced, %d new, %d updated, %d errors",
			result.SyncedCount, result.NewCount, result.UpdatedCount, result.ErrorCount)
	}

	for {
		select {
		case <-ticker.C:
			result, err := syncSvc.SyncAll()
			if err != nil {
				log.Printf("Sync failed: %v", err)
			} else {
				log.Printf("Sync: %d synced, %d new, %d updated, %d errors",
					result.SyncedCount, result.NewCount, result.UpdatedCount, result.ErrorCount)
			}
		case <-ctx.Done():
			return
		}
	}
}

// handleCommand 处理命令行命令
func handleCommand(args []string) {
	command := args[0]

	switch command {
	case "test":
		testCookie()
	case "init":
		if len(args) < 4 {
			log.Println("Error: init command requires TOKEN, CSRF_TOKEN, and UID")
			log.Println("Usage: copytrading-sync-worker init <token> <csrf_token> <uid>")
			os.Exit(1)
		}
		initCookie(args[1], args[2], args[3])
	case "sync":
		runManualSync()
	case "help", "-h", "--help":
		printUsage()
	default:
		log.Printf("Unknown command: %s", command)
		printUsage()
		os.Exit(1)
	}
}

// testCookie 测试 Cookie 有效性
func testCookie() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}
	cfg := config.C

	db, err := repository.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	repo := repository.NewCopyTradingRepository(db.DB)
	cookieMgr := copytrading.NewCookieManager(cfg.CopyTrading.EncryptionKey)
	webClient := copytrading.NewGateWebClient(cfg.CopyTrading.UserAgent)

	cookie, err := repo.GetActiveCookie()
	if err != nil {
		log.Fatalf("No active cookie found: %v", err)
	}

	gateCookie, err := cookieMgr.DecryptGateCookie(
		cookie.Token,
		cookie.CsrfToken,
		cookie.Uid,
		cookie.ExpiresAt,
	)
	if err != nil {
		log.Fatalf("Failed to decrypt cookie: %v", err)
	}

	log.Printf("Cookie Status:")
	log.Printf("  UID: %s", gateCookie.Uid)
	log.Printf("  Expires: %s", gateCookie.ExpiresAt.Format(time.RFC3339))
	log.Printf("  Is Expired: %v", cookieMgr.IsExpired(gateCookie))
	log.Printf("  Time until expiry: %s", formatExpiry(gateCookie.ExpiresAt))

	// 测试 API 调用
	params := copytrading.DefaultTraderListParams()
	params.Page = 1
	params.PageSize = 5

	result, err := webClient.GetTraderList(gateCookie, params)
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}

	log.Printf("API Test Success!")
	log.Printf("  Total traders: %d", result.Data.TotalCount)
	log.Printf("  Returned: %d", len(result.Data.List))
	printTraders(result.Data.List, 5)
}

// initCookie 初始化 Cookie
func initCookie(token, csrfToken, uid string) {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}
	cfg := config.C

	db, err := repository.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Println("Running database migration...")
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	repo := repository.NewCopyTradingRepository(db.DB)
	cookieMgr := copytrading.NewCookieManager(cfg.CopyTrading.EncryptionKey)

	gateCookie, err := cookieMgr.NewGateCookie(token, csrfToken, uid)
	if err != nil {
		log.Fatalf("Failed to parse cookie: %v", err)
	}

	encryptedToken, encryptedCsrf, err := cookieMgr.EncryptGateCookie(gateCookie)
	if err != nil {
		log.Fatalf("Failed to encrypt cookie: %v", err)
	}

	cookie := &models.PlatformCookie{
		Token:        encryptedToken,
		CsrfToken:    encryptedCsrf,
		Uid:          uid,
		ExpiresAt:    &gateCookie.ExpiresAt,
		Status:       "active",
		LastSyncedAt: nil,
	}

	if err := repo.CreateOrUpdateCookie(cookie); err != nil {
		log.Fatalf("Failed to save cookie: %v", err)
	}

	log.Println("Cookie saved successfully!")
	log.Printf("  UID: %s", uid)
	log.Printf("  Expires: %s", gateCookie.ExpiresAt.Format(time.RFC3339))
	log.Printf("  Days until expiry: %.1f", time.Until(gateCookie.ExpiresAt).Hours()/24)
}

// runManualSync 手动触发同步
func runManualSync() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}
	cfg := config.C

	db, err := repository.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	repo := repository.NewCopyTradingRepository(db.DB)
	cookieMgr := copytrading.NewCookieManager(cfg.CopyTrading.EncryptionKey)
	webClient := copytrading.NewGateWebClient(cfg.CopyTrading.UserAgent)
	syncSvc := copytrading.NewSyncService(repo, cookieMgr, webClient)

	result, err := syncSvc.SyncAll()
	if err != nil {
		log.Printf("Manual sync failed: %v", err)
		os.Exit(1)
	}

	log.Printf("Manual sync completed: %s", result.String())

	var existingTraders int64
	db.DB.Model(&models.CopyTrader{}).Count(&existingTraders)
	log.Printf("Total traders in database: %d", existingTraders)
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println(`
Copy Trading Sync Worker

Usage:
  copytrading-sync-worker [command]

Commands:
  (none)       运行定时同步任务
  test         测试 Cookie 有效性
  init         初始化 Cookie (需要 TOKEN, CSRF, UID 参数)
  sync         手动触发一次同步
  help         显示帮助信息

Examples:
  # 运行定时同步
  copytrading-sync-worker

  # 测试 Cookie
  copytrading-sync-worker test

  # 初始化 Cookie
  copytrading-sync-worker init <token> <csrf_token> <uid>

  # 手动同步
  copytrading-sync-worker sync

Environment Variables:
  COPYTRADING_SYNC_INTERVAL    同步间隔（分钟），默认 30
  ENCRYPTION_KEY               加密密钥（32 字节）
  COPYTRADING_USER_AGENT       User-Agent 字符串
  DATABASE_URL                 数据库连接字符串
`)
}

// formatExpiry 格式化过期时间
func formatExpiry(expiresAt time.Time) string {
	if time.Now().After(expiresAt) {
		return "EXPIRED"
	}
	hours := time.Until(expiresAt).Hours()
	if hours < 24 {
		return fmt.Sprintf("%.1f hours", hours)
	}
	return fmt.Sprintf("%.1f days", hours/24)
}

// printTraders 打印交易员列表
func printTraders(traders []copytrading.TraderInfo, max int) {
	count := len(traders)
	if max > 0 && count > max {
		count = max
	}

	fmt.Printf("\n%-4s %-20s %-12s %-15s %-10s %-10s\n", "#", "Name", "Profit", "ROI (%)", "Win Rate", "Followers")
	fmt.Println(strings.Repeat("-", 90))

	for i := 0; i < count; i++ {
		t := traders[i]
		name := t.UserInfo.Nickname
		if name == "" {
			name = t.UserInfo.HideName
		}
		fmt.Printf("%-4d %-20s $%-14.2f %-10.2f %-10.2f %-10d\n",
			i+1, name, float64(t.Profit), float64(t.ProfitRate)*100, float64(t.WinRate)*100, t.CurrFollowNum)
	}

	if len(traders) > max {
		fmt.Printf("... and %d more\n", len(traders)-max)
	}
}
