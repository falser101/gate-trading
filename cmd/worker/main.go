package main

import (
	"context"
	"log"
	"time"

	"github.com/falser101/gate-trading/internal/config"
	"github.com/falser101/gate-trading/internal/gateway"
	"github.com/falser101/gate-trading/internal/models"
	"github.com/falser101/gate-trading/internal/repository"
	"github.com/falser101/gate-trading/internal/service/strategy"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	// 初始化数据库
	db, err := repository.NewDatabase(config.C.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 初始化服务
	strategyRepo := repository.NewStrategyRepository(db.DB)
	orderRepo := repository.NewOrderRepository(db.DB)
	logRepo := repository.NewStrategyLogRepository(db.DB)
	gateFactory := gateway.NewGateClientFactory(config.C.Gate.BaseURL)

	strategyEngine := strategy.NewStrategyEngine(strategyRepo, orderRepo, logRepo, gateFactory)

	log.Println("Worker started, loading strategies...")

	// 加载所有运行中的策略
	strategies, err := strategyRepo.GetByStatus(strategy.StatusRunning)
	if err != nil {
		log.Printf("Failed to load strategies: %v", err)
	}

	for _, s := range strategies {
		// 获取用户信息
		var user models.User
		if err := db.DB.Model(&s).Association("User").Find(&user); err != nil {
			log.Printf("Failed to get user for strategy %d: %v", s.ID, err)
			continue
		}

		if err := strategyEngine.LoadStrategy(&s, &user); err != nil {
			log.Printf("Failed to load strategy %d: %v", s.ID, err)
			continue
		}
		log.Printf("Loaded strategy %d (%s - %s)", s.ID, s.Type, s.Symbol)
	}

	// 主循环：每 10 秒检查一次行情并执行策略
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for range ticker.C {
		// 获取所有运行中的策略
		strategies, err := strategyRepo.GetByStatus(strategy.StatusRunning)
		if err != nil {
			log.Printf("Failed to get strategies: %v", err)
			continue
		}

		for _, s := range strategies {
			// 获取用户 API Key
			var user models.User
			if err := db.DB.Model(&s).Association("User").Find(&user); err != nil {
				continue
			}

			// 获取当前价格
			client := gateFactory.GetClient(user.GateAPIKey, user.GateAPISecret)
			tickerInfo, err := client.GetTicker(s.Symbol)
			if err != nil {
				log.Printf("Failed to get ticker for %s: %v", s.Symbol, err)
				continue
			}

			price := tickerInfo.Last

			// 执行策略
			strat := strategyEngine.GetStrategy(s.ID)
			if strat != nil {
				if err := strat.OnTick(ctx, price); err != nil {
					log.Printf("Strategy %d OnTick error: %v", s.ID, err)
				}
			}
		}
	}
}
