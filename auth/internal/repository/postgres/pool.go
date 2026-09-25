package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Pool struct {
	*gorm.DB
}

func NewPool(
	cfg config,
	domains ...any,
) (*Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DefaultContextTimeout: cfg.Timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(domains...); err != nil {
		return nil, fmt.Errorf("failed to migrate domain: %w", err)
	}

	return &Pool{
		DB: db,
	}, nil
}
