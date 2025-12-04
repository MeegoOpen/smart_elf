package database

import (
	"log"
	"smart_elf_standalone/internal/model"
	"smart_elf_standalone/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 配置GORM日志
	logLevel := logger.Info
	if !cfg.Debug {
		logLevel = logger.Warn
	}

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var db *gorm.DB
	var err error

	// 根据DSN判断使用MySQL还是SQLite
	if cfg.DSN == "sqlite://./smart_elf.db" {
		// 使用SQLite
		log.Printf("信息: 使用SQLite数据库: ./smart_elf.db")
		db, err = gorm.Open(sqlite.Open("./smart_elf.db"), gormConfig)
	} else {
		// 使用MySQL
		log.Printf("信息: 使用MySQL数据库")
		db, err = gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	}

	if err != nil {
		log.Printf("错误: 连接数据库失败: %v", err)
		return nil, err
	}

	// 配置连接池（对两种数据库都适用）
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("错误: 获取数据库连接池失败: %v", err)
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	log.Printf("信息: 数据库连接池配置完成: max_idle_conns=%d, max_open_conns=%d, debug=%v", cfg.MaxIdleConns, cfg.MaxOpenConns, cfg.Debug)

	return db, nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	log.Printf("信息: 开始数据库迁移")

	// 迁移模型
	err := db.AutoMigrate(
		&model.AppConfig{},
	)

	if err != nil {
		log.Printf("错误: 数据库迁移失败: %v", err)
		return err
	}

	log.Printf("信息: 数据库迁移成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		log.Printf("错误: 关闭数据库连接失败: %v", err)
		return err
	}

	log.Printf("信息: 数据库连接关闭成功")
	return nil
}
