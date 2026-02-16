package database

import (
	"context"
	"fmt"
	"orm"

	"github.com/rs/zerolog/log"

	"pocket-shop/config"
)

func ProvideDatabase(cfg *config.Config) (*orm.Client, error) {
	dsn := cfg.PostgresDSN()
	client, err := orm.NewPostgresClient(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Msg("Connected to PostgreSQL")

	return client, nil
}

func EnsureSchema(ctx context.Context, client *orm.Client) error {
	return client.Schema.Create(ctx)
}
