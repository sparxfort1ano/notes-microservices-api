package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewDatabase(cfg config) (*mongo.Database, error) {
	dsn := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?tls=%t",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.SSLMode,
	)

	client, err := mongo.Connect(options.Client().ApplyURI(dsn).SetTimeout(cfg.Timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongo database: %w", err)
	}

	return client.Database(cfg.DB), nil
}

func CloseDB(db *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), NewConfigMust().Timeout)
	defer cancel()

	return db.Disconnect(ctx)
}
