package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"ecommerce/user-service/internal/model"
	"ecommerce/user-service/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var sqlitePath = getDefaultSQLitePath()

func getDefaultSQLitePath() string {
	//从环境变量获取，如果没有则使用默认路径
	if path := os.Getenv("SQLITE_PATH"); path != "" {
		return path
	}
	//默认路径为当前工作目录下的data/user.db
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return filepath.Join(cwd, "data", "user.db")
}

// 建立数据库
func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, string, error) {
	//先尝试MySQL
	db, err := tryMySQL(cfg)
	if err == nil {
		if err := autoMigrate(db); err != nil {
			log.Printf("MySQL 表迁移失败: %v", err)
		}
		return db, "mysql", nil
	}
	log.Printf("Failed to connect to MySQL: %v, falling back to SQLite", err)
	//降级到SQLite
	db, err = trySQLite()
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to both MySQL and SQLite: %v", err)
	}
	if err := autoMigrate(db); err != nil {
		log.Printf("SQLite 表迁移失败: %v", err)
	}

	return db, "sqlite", nil
}

// mysql连接
func tryMySQL(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	//设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// sqlite连接
func trySQLite() (*gorm.DB, error) {
	dir := filepath.Dir(sqlitePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for SQLite: %v", err)
	}
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	return db, nil
}

// 自动建表
func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return fmt.Errorf("迁移User表失败: %v", err)
	}
	log.Println("数据库表迁移成功")
	return nil
}
