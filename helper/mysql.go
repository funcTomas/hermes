package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/funcTomas/hermes/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlClient struct {
	cfg *config.MysqlConfig
}

func NewMysqlClient(cfg *config.MysqlConfig) *MysqlClient {
	return &MysqlClient{cfg: cfg}
}

func (c *MysqlClient) Connect(ctx context.Context) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.cfg.User, c.cfg.Password, c.cfg.Host, c.cfg.Port, c.cfg.Database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(c.cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.cfg.MaxOpenConns)
	if lifetime, parseErr := time.ParseDuration(c.cfg.ConnMaxLifetime); parseErr == nil {
		sqlDB.SetConnMaxLifetime(lifetime)
	} else {
		log.Printf("Warning: Invalid ConnMaxLifetime format: %s, using default", c.cfg.ConnMaxLifetime)
	}

	log.Println("Connected to MySQL with GORM")
	return db, nil
}

func (c *MysqlClient) Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
