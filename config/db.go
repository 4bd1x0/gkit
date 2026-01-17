package config

import (
	"database/sql"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstances = make(map[string]*gorm.DB)
	dbMutex     sync.RWMutex
)

func NewDB(name string) (*gorm.DB, error) {
	dbMutex.RLock()
	if db, exists := dbInstances[name]; exists {
		dbMutex.RUnlock()
		return db, nil
	}
	dbMutex.RUnlock()

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db, exists := dbInstances[name]; exists {
		return db, nil
	}

	cfg, exists := config.DB[name]
	if !exists {
		return nil, fmt.Errorf("database config '%s' not found", name)
	}

	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		if err := conn.Ping(); err != nil {
			return nil, err
		}
		dialector = mysql.New(mysql.Config{Conn: conn})

	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
		dialector = postgres.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}

	dbInstances[name] = db
	return db, nil
}

func DefaultDB() (*gorm.DB, error) {
	return NewDB("primary")
}
