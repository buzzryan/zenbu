package rdbutil

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/buzzryan/zenbu/internal/config"
)

func MustConnectMySQL(cfg config.MySQLConfig) *gorm.DB {
	gormLogger := logger.Discard
	if cfg.LogEnabled {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		// SkipDefaultTransaction disables the default transaction for gorm operations.
		// It gives 30%+ performance improvement.
		SkipDefaultTransaction: true,
		// Logger is the logger used for this gorm.DB instance.
		Logger: gormLogger,
		// TranslateError translates the error returned by the driver to gorm's error.
		// For example, it will translate the "record not found" error to gorm.ErrRecordNotFound.
		TranslateError: true,
	})
	if err != nil {
		log.Panicf("failed to connect mysql(DSN: %v): %v\n", cfg.DSN(), err)
	}

	return setConnectionPool(db)
}

// setConnectionPool sets the connection pool settings for the given gorm.DB.
func setConnectionPool(db *gorm.DB) *gorm.DB {
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(70)
	sqlDB.SetMaxOpenConns(75)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db
}
