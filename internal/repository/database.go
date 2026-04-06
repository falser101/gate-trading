package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/falser/gate-trading/internal/config"
	"github.com/falser/gate-trading/internal/models"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg config.DBConfig) (*Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

// 自动迁移表结构
func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&models.User{},
		&models.Strategy{},
		&models.Order{},
		&models.Trade{},
		&models.StrategyLog{},
	)
}
