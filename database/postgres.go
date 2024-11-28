package database

import (
	"crowdfunding/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(cfg config.DBConfig) (db *gorm.DB, err error) {
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Dapatkan *sql.DB dari GORM untuk mengatur connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Atur connection pool
	sqlDB.SetMaxIdleConns(int(cfg.ConnectionPool.MaxIdleConnection))
	sqlDB.SetMaxOpenConns(int(cfg.ConnectionPool.MaxOpenConnection))
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnectionPool.MaxIdletimeConnection) * time.Second)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnectionPool.MaxLifetimeConnection) * time.Second)

	log.Println("Connected to PostgreSQL with GORM and connection pool configured!")
	return db, nil
}

// func ConnectPostgres(cfg config.DBConfig) (db *sqlx.DB, err error) {
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		cfg.Host,
// 		cfg.Port,
// 		cfg.User,
// 		cfg.Password,
// 		cfg.Name,
// 	)

// 	db, err = sqlx.Open("postgres", dsn)
// 	if err != nil {
// 		return
// 	}

// 	if err = db.Ping(); err != nil {
// 		db = nil
// 		return
// 	}

// 	db.SetConnMaxIdleTime(time.Duration(cfg.ConnectionPool.MaxIdletimeConnection) * time.Second)
// 	db.SetConnMaxLifetime(time.Duration(cfg.ConnectionPool.MaxLifetimeConnection) * time.Second)
// 	db.SetMaxOpenConns(int(cfg.ConnectionPool.MaxOpenConnetcion))
// 	db.SetMaxIdleConns(int(cfg.ConnectionPool.MaxIdleConnection))

// 	return
// }
