package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/config"
)

type courierRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewCourierRepositoryPostgres(ctx context.Context, cfg config.Config) (*courierRepositoryPostgres, error) {
	const op = "repository.postgres.NewCourierRepositoryPostgres"

	pgxCfg, err := parseConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &courierRepositoryPostgres{pool: conn}, nil
}

func parseConfig(cfg config.Config) (*pgxpool.Config, error) {
	const op = "repository.postgres.parseConfig"

	storageDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresStorage.User,
		cfg.PostgresStorage.Password,
		cfg.PostgresStorage.Host,
		cfg.PostgresStorage.DbPort,
		cfg.PostgresStorage.DbName,
	)

	pgxCfg, err := pgxpool.ParseConfig(storageDSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns
	pgxCfg.MaxConnLifetime = cfg.MaxConnLifeTime

	return pgxCfg, nil
}
