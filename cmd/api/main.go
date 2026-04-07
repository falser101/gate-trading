package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"

	"github.com/falser101/gate-trading/internal/config"
	"github.com/falser101/gate-trading/internal/gateway"
	"github.com/falser101/gate-trading/internal/middleware"
	"github.com/falser101/gate-trading/internal/repository"
	"github.com/falser101/gate-trading/internal/service/auth"
	"github.com/falser101/gate-trading/internal/service/copytrading"
	"github.com/falser101/gate-trading/internal/service/strategy"
	gateapi "github.com/gate/gateapi-go/v7"
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

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// 初始化服务
	userRepo := repository.NewUserRepository(db.DB)
	strategyRepo := repository.NewStrategyRepository(db.DB)
	orderRepo := repository.NewOrderRepository(db.DB)
	logRepo := repository.NewStrategyLogRepository(db.DB)
	copyTradingRepo := repository.NewCopyTradingRepository(db.DB)

	authService := auth.NewAuthService(userRepo, config.C.AppSecret)
	gateFactory := gateway.NewGateClientFactory(config.C.Gate.BaseURL)
	strategyEngine := strategy.NewStrategyEngine(strategyRepo, orderRepo, logRepo, gateFactory)
	cookieMgr := copytrading.NewCookieManager(config.C.CopyTrading.EncryptionKey)
	webClient := copytrading.NewGateWebClient(config.C.CopyTrading.UserAgent)
	copyTradingSvc := copytrading.NewSyncService(copyTradingRepo, cookieMgr, webClient)

	// 设置 Gin 模式
	ginMode := gin.DebugMode
	if config.C.Env == "release" {
		ginMode = gin.ReleaseMode
	} else if config.C.Env == "test" {
		ginMode = gin.TestMode
	}
	gin.SetMode(ginMode)
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 认证路由
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", handleRegister(authService))
		authGroup.POST("/login", handleLogin(authService))
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(authService))
	{
		// 用户信息
		protected.GET("/user", handleGetUser(userRepo))

		// 绑定 API Key
		protected.POST("/user/api-key", handleBindAPIKey(authService))

		// 策略管理
		strategies := protected.Group("/strategies")
		{
			strategies.POST("", handleCreateStrategy(strategyEngine))
			strategies.GET("", handleListStrategies(strategyRepo))
			strategies.GET("/:id", handleGetStrategy(strategyRepo))
			strategies.POST("/:id/start", handleStartStrategy(strategyEngine))
			strategies.POST("/:id/stop", handleStopStrategy(strategyEngine))
			strategies.DELETE("/:id", handleDeleteStrategy(strategyRepo))
		}

		// 订单查询
		protected.GET("/orders", handleListOrders(orderRepo))

		// 行情
		market := protected.Group("/market")
		{
			market.GET("/ticker/:symbol", handleGetTicker(gateFactory))
		}

		// 账户
		protected.GET("/account/balance", handleGetBalance(userRepo, gateFactory))
		protected.GET("/account/detail", handleGetAccountDetail(userRepo, gateFactory))

		// 合约交易
		futures := protected.Group("/futures")
		{
			futures.GET("/accounts", handleGetFuturesAccount(userRepo, gateFactory))
			futures.GET("/positions", handleGetFuturesPositions(userRepo, gateFactory))
			futures.GET("/positions/:contract", handleGetFuturesPosition(userRepo, gateFactory))
			futures.POST("/orders", handleCreateFuturesOrder(userRepo, gateFactory))
			futures.GET("/orders", handleGetFuturesOrders(userRepo, gateFactory))
			futures.DELETE("/orders/:id", handleCancelFuturesOrder(userRepo, gateFactory))
			futures.POST("/positions/:contract/leverage", handleSetFuturesLeverage(userRepo, gateFactory))
			futures.POST("/positions/close", handleCloseFuturesPosition(userRepo, gateFactory))
			futures.GET("/tickers", handleGetFuturesTickers(gateFactory))
			futures.GET("/tickers/:contract", handleGetFuturesTicker(gateFactory))
			futures.GET("/contracts", handleGetFuturesContracts(gateFactory))
		}

		// 跟单交易（公开接口，所有用户可访问）
		api := r.Group("/api")
		{
			api.GET("/copytrading/traders", handleGetCopyTraders(copyTradingSvc))
			api.GET("/copytrading/traders/:id", handleGetCopyTraderDetail(copyTradingSvc))
			api.GET("/copytrading/traders/:id/stats", handleGetCopyTraderStats(copyTradingSvc))
		}
	}

	// 管理员接口（暂时简单处理，后续可添加权限控制）
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(authService))
	{
		// Cookie 管理
		admin.POST("/copytrading/cookie", handleSaveCookie(copyTradingSvc, cookieMgr, copyTradingRepo))
		admin.GET("/copytrading/cookie", handleGetCookie(copyTradingRepo))
		admin.DELETE("/copytrading/cookie", handleDeleteCookie(copyTradingRepo))
		admin.POST("/copytrading/sync", handleManualSync(copyTradingSvc))
		// 通知
		admin.GET("/notifications", handleGetNotifications(copyTradingRepo))
		admin.POST("/notifications/mark-read", handleMarkNotificationRead(copyTradingRepo))
	}

	// 启动服务
	addr := ":" + strconv.Itoa(config.C.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRegister(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := authService.Register(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func handleLogin(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := authService.Login(&req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func handleGetUser(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":          user.ID,
			"email":       user.Email,
			"api_key_set": user.GateAPIKey != "",
		})
	}
}

func handleBindAPIKey(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var req struct {
			APIKey    string `json:"api_key"`
			APISecret string `json:"api_secret"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := authService.BindAPIKey(userID, req.APIKey, req.APISecret); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "API key bound successfully"})
	}
}

func handleCreateStrategy(engine *strategy.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		var req struct {
			Type   string         `json:"type"`
			Symbol string         `json:"symbol"`
			Config map[string]any `json:"config"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		strat, err := engine.CreateStrategy(userID, req.Type, req.Symbol, req.Config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, strat)
	}
}

func handleListStrategies(repo *repository.StrategyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		strategies, err := repo.GetByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": strategies})
	}
}

func handleGetStrategy(repo *repository.StrategyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		// Simple implementation, parse id as needed
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "get strategy"})
	}
}

func handleStartStrategy(engine *strategy.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		// Simple implementation, parse id as needed
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "start strategy"})
	}
}

func handleStopStrategy(engine *strategy.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		// Simple implementation, parse id as needed
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "stop strategy"})
	}
}

func handleDeleteStrategy(repo *repository.StrategyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		// Simple implementation, parse id as needed
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "delete strategy"})
	}
}

func handleListOrders(repo *repository.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []any{}})
	}
}

func handleGetTicker(factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		client := factory.GetClient("", "")

		// 使用官方 SDK 获取 ticker
		tickers, _, err := client.SpotApi.ListTickers(context.Background(), &gateapi.ListTickersOpts{CurrencyPair: optional.NewString(symbol)})
		if err != nil || len(tickers) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": tickers[0]})
	}
}

func handleGetBalance(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		// 从数据库获取用户的 API Key
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		balance, _, err := client.WalletApi.GetTotalBalance(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 转换为自定义响应格式
		response := map[string]any{
			"total": map[string]any{
				"amount":         balance.Total.Amount,
				"currency":       balance.Total.Currency,
				"unrealised_pnl": balance.Total.UnrealisedPnl,
				"borrowed":       balance.Total.Borrowed,
			},
			"details": map[string]any{},
		}

		// 转换 details
		detailsMap := make(map[string]any)
		for k, v := range balance.Details {
			detailsMap[k] = map[string]any{
				"amount":         v.Amount,
				"currency":       v.Currency,
				"unrealised_pnl": v.UnrealisedPnl,
				"borrowed":       v.Borrowed,
			}
		}
		response["details"] = detailsMap

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetAccountDetail(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		// 从数据库获取用户的 API Key
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		detail, _, err := client.AccountApi.GetAccountDetail(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 转换为自定义响应格式
		response := map[string]any{
			"user_id": detail.UserId,
			"tier":    detail.Tier,
		}

		// 转换 key 信息 - detail.Key 是值类型不是指针
		response["key"] = map[string]any{
			"mode": detail.Key.Mode,
		}

		// 转换 IP 白名单
		if len(detail.IpWhitelist) > 0 {
			response["ip_whitelist"] = detail.IpWhitelist
		}

		// 转换允许的交易对
		if len(detail.CurrencyPairs) > 0 {
			response["currency_pairs"] = detail.CurrencyPairs
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

// ParseDecimal helper - simple implementation
func parseDecimal(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

// ========== Futures Handlers ==========

func handleGetFuturesAccount(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		account, _, err := client.FuturesApi.ListFuturesAccounts(ctx, "usdt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"total":              account.Total,
			"unrealised_pnl":     account.UnrealisedPnl,
			"available":          account.Available,
			"order_margin":       account.OrderMargin,
			"position_margin":    account.PositionMargin,
			"maintenance_margin": account.MaintenanceMargin,
			"currency":           account.Currency,
			"in_dual_mode":       account.InDualMode,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetFuturesPositions(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		positions, _, err := client.FuturesApi.ListPositions(ctx, "usdt", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []map[string]interface{}
		for _, pos := range positions {
			response = append(response, map[string]interface{}{
				"contract":          pos.Contract,
				"size":              pos.Size,
				"leverage":          pos.Leverage,
				"entry_price":       pos.EntryPrice,
				"mark_price":        pos.MarkPrice,
				"liq_price":         pos.LiqPrice,
				"margin":            pos.Margin,
				"unrealised_pnl":    pos.UnrealisedPnl,
				"realised_pnl":      pos.RealisedPnl,
				"initial_margin":    pos.InitialMargin,
				"maintenance_margin": pos.MaintenanceMargin,
				"adl_ranking":       pos.AdlRanking,
				"mode":              pos.Mode,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetFuturesPosition(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		contract := c.Param("contract")
		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		position, _, err := client.FuturesApi.GetPosition(ctx, "usdt", contract)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"contract":           position.Contract,
			"size":               position.Size,
			"leverage":           position.Leverage,
			"entry_price":        position.EntryPrice,
			"mark_price":         position.MarkPrice,
			"liq_price":          position.LiqPrice,
			"margin":             position.Margin,
			"unrealised_pnl":     position.UnrealisedPnl,
			"realised_pnl":       position.RealisedPnl,
			"initial_margin":     position.InitialMargin,
			"maintenance_margin": position.MaintenanceMargin,
			"adl_ranking":        position.AdlRanking,
			"mode":               position.Mode,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleCreateFuturesOrder(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		var req struct {
			Contract   string `json:"contract"`
			Size       string `json:"size"`
			Price      string `json:"price"`
			Tif        string `json:"tif"`
			ReduceOnly bool   `json:"reduce_only"`
			Close      bool   `json:"close"`
			Text       string `json:"text"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		order := gateapi.FuturesOrder{
			Contract:   req.Contract,
			Size:       req.Size,
			Price:      req.Price,
			Tif:        req.Tif,
			ReduceOnly: req.ReduceOnly,
			Close:      req.Close,
			Text:       req.Text,
		}

		createdOrder, _, err := client.FuturesApi.CreateFuturesOrder(ctx, "usdt", order, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"id":          createdOrder.Id,
			"contract":    createdOrder.Contract,
			"size":        createdOrder.Size,
			"price":       createdOrder.Price,
			"status":      createdOrder.Status,
			"left":        createdOrder.Left,
			"fill_price":  createdOrder.FillPrice,
			"create_time": createdOrder.CreateTime,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetFuturesOrders(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		status := c.Query("status")
		if status == "" {
			status = "open"
		}
		limitStr := c.Query("limit")
		limit := int32(100)
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = int32(l)
			}
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		orders, _, err := client.FuturesApi.ListFuturesOrders(ctx, "usdt", status, &gateapi.ListFuturesOrdersOpts{
			Limit: optional.NewInt32(limit),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []map[string]interface{}
		for _, order := range orders {
			response = append(response, map[string]interface{}{
				"id":          order.Id,
				"contract":    order.Contract,
				"size":        order.Size,
				"price":       order.Price,
				"status":      order.Status,
				"left":        order.Left,
				"fill_price":  order.FillPrice,
				"create_time": order.CreateTime,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleCancelFuturesOrder(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		orderId := c.Param("id")
		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		cancelledOrder, _, err := client.FuturesApi.CancelFuturesOrder(ctx, "usdt", orderId, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"id":          cancelledOrder.Id,
			"contract":    cancelledOrder.Contract,
			"status":      cancelledOrder.Status,
			"finish_as":   cancelledOrder.FinishAs,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleSetFuturesLeverage(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		contract := c.Param("contract")
		var req struct {
			Leverage   string `json:"leverage"`
			MarginMode string `json:"margin_mode"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		position, _, err := client.FuturesApi.UpdateContractPositionLeverage(ctx, "usdt", contract, req.Leverage, req.MarginMode, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"contract":  position.Contract,
			"leverage":  position.Leverage,
			"mode":      position.Mode,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleCloseFuturesPosition(userRepo *repository.UserRepository, factory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if user.GateAPIKey == "" || user.GateAPISecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先配置 Gate API Key", "api_key_set": false})
			return
		}

		var req struct {
			Contract string `json:"contract"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		client := factory.GetClient(user.GateAPIKey, user.GateAPISecret)
		ctx := factory.GetContext(user.GateAPIKey, user.GateAPISecret)

		order, _, err := client.FuturesApi.CreateFuturesOrder(ctx, "usdt", gateapi.FuturesOrder{
			Contract: req.Contract,
			Size:     "0",
			Price:    "0",
			Close:    true,
			Tif:      "ioc",
		}, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"id":          order.Id,
			"contract":    order.Contract,
			"size":        order.Size,
			"fill_price":  order.FillPrice,
			"finish_as":   order.FinishAs,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetFuturesTickers(gateFactory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := gateFactory.GetClient("", "")
		ctx := context.Background()

		tickers, _, err := client.FuturesApi.ListFuturesTickers(ctx, "usdt", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []map[string]interface{}
		for _, ticker := range tickers {
			response = append(response, map[string]interface{}{
				"contract":          ticker.Contract,
				"last":              ticker.Last,
				"change_percentage": ticker.ChangePercentage,
				"volume_24h":        ticker.Volume24h,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetFuturesTicker(gateFactory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		contract := c.Param("contract")
		client := gateFactory.GetClient("", "")
		ctx := context.Background()

		// 使用 ListFuturesTickers 然后过滤
		tickers, _, err := client.FuturesApi.ListFuturesTickers(ctx, "usdt", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var ticker *gateapi.FuturesTicker
		for _, t := range tickers {
			if t.Contract == contract {
				ticker = &t
				break
			}
		}

		if ticker == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
			return
		}

		response := map[string]interface{}{
			"contract":          ticker.Contract,
			"last":              ticker.Last,
			"mark_price":        ticker.MarkPrice,
			"index_price":       ticker.IndexPrice,
			"change_percentage": ticker.ChangePercentage,
			"volume_24h":        ticker.Volume24h,
			"volume_24h_usd":    ticker.Volume24hUsd,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func handleGetFuturesContracts(gateFactory *gateway.GateClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := gateFactory.GetClient("", "")
		ctx := context.Background()

		contracts, _, err := client.FuturesApi.ListFuturesContracts(ctx, "usdt", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []map[string]interface{}
		for _, contract := range contracts {
			response = append(response, map[string]interface{}{
				"name":              contract.Name,
				"type":              contract.Type,
				"quanto_multiplier": contract.QuantoMultiplier,
				"leverage_min":      contract.LeverageMin,
				"leverage_max":      contract.LeverageMax,
				"mark_price":        contract.MarkPrice,
				"index_price":       contract.IndexPrice,
				"funding_rate":      contract.FundingRate,
				"order_size_min":    contract.OrderSizeMin,
				"order_size_max":    contract.OrderSizeMax,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

// ==================== Copy Trading Handlers ====================
