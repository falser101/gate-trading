package repository

import (
	"github.com/falser/gate-trading/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) UpdateAPIKey(userID uint, apiKey, apiSecret string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"gate_api_key":    apiKey,
			"gate_api_secret": apiSecret,
		}).Error
}

type StrategyRepository struct {
	db *gorm.DB
}

func NewStrategyRepository(db *gorm.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

func (r *StrategyRepository) Create(strategy *models.Strategy) error {
	return r.db.Create(strategy).Error
}

func (r *StrategyRepository) GetByID(id uint) (*models.Strategy, error) {
	var strategy models.Strategy
	err := r.db.Preload("Orders").First(&strategy, id).Error
	return &strategy, err
}

func (r *StrategyRepository) GetByUserID(userID uint) ([]models.Strategy, error) {
	var strategies []models.Strategy
	err := r.db.Where("user_id = ?", userID).Find(&strategies).Error
	return strategies, err
}

func (r *StrategyRepository) GetByStatus(status string) ([]models.Strategy, error) {
	var strategies []models.Strategy
	err := r.db.Where("status = ?", status).Find(&strategies).Error
	return strategies, err
}

func (r *StrategyRepository) Update(strategy *models.Strategy) error {
	return r.db.Save(strategy).Error
}

func (r *StrategyRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Strategy{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *StrategyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Strategy{}, id).Error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) GetByStrategyID(strategyID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("strategy_id = ?", strategyID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetByGateOrderID(gateOrderID string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("gate_order_id = ?", gateOrderID).First(&order).Error
	return &order, err
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) UpdateStatus(gateOrderID, status string) error {
	return r.db.Model(&models.Order{}).
		Where("gate_order_id = ?", gateOrderID).
		Update("status", status).Error
}

type TradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

func (r *TradeRepository) Create(trade *models.Trade) error {
	return r.db.Create(trade).Error
}

func (r *TradeRepository) GetByOrderID(orderID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.Where("order_id = ?", orderID).Find(&trades).Error
	return trades, err
}

type StrategyLogRepository struct {
	db *gorm.DB
}

func NewStrategyLogRepository(db *gorm.DB) *StrategyLogRepository {
	return &StrategyLogRepository{db: db}
}

func (r *StrategyLogRepository) Create(log *models.StrategyLog) error {
	return r.db.Create(log).Error
}

func (r *StrategyLogRepository) GetByStrategyID(strategyID uint, limit int) ([]models.StrategyLog, error) {
	var logs []models.StrategyLog
	err := r.db.Where("strategy_id = ?", strategyID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
