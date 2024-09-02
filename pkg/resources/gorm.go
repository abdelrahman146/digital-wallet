package resources

import (
	"digital-wallet/pkg/config"
	"digital-wallet/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := config.GetConfig().GetDbConnectionString()
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		//Logger:         gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		logger.GetLogger().Panic("failed to connect to database", logger.Field("error", err))
	}
	db, err := gormDB.DB()
	if err != nil {
		logger.GetLogger().Panic("failed to get database connection", logger.Field("error", err))
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	return gormDB
}

func CloseDB(gormDB *gorm.DB) {
	db, err := gormDB.DB()
	if err != nil {
		logger.GetLogger().Panic("failed to get database connection", logger.Field("error", err))
	}
	err = db.Close()
	if err != nil {
		logger.GetLogger().Panic("failed to close database connection", logger.Field("error", err))
	}
}
