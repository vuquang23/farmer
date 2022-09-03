package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	cfg "farmer/internal/pkg/config"
	"farmer/pkg/errors"
)

var db *gorm.DB

func InitDB() error {
	if db != nil {
		return nil
	}

	dbCfg := cfg.Instance().DB
	gormDB, err := gorm.Open(
		mysql.Open(dbCfg.DSN()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(dbCfg.LogLevel)),
		},
	)
	if err != nil {
		return errors.NewInfraErrorDBConnect(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return errors.NewInfraErrorDBConnect(err)
	}

	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.ConnLifeTime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		return errors.NewInfraErrorDBConnect(err)
	}

	db = gormDB
	return nil
}

func Instance() *gorm.DB {
	return db
}
