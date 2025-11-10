package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/config"
	"time"
)

const (
	MAX_RETRIES = 5
	BASE_DELAY  = 100 * time.Millisecond
)

type courierRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewCourierRepositoryPostgres(ctx context.Context, cfg config.PostgresStorage) (*courierRepositoryPostgres, error) {
	const op = "repository.postgres.NewCourierRepositoryPostgres"

	pgxCfg, err := parseConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := 0; i <= MAX_RETRIES; i++ {
		err := conn.Ping(ctx)

		if err != nil && i == MAX_RETRIES {
			return nil, fmt.Errorf("%s: %w", op, err)
		} else if err != nil {
			//exponential backoff
			backoff := time.Duration(1<<i) * BASE_DELAY
			time.Sleep(backoff)
		} else {
			break
		}
	}
	return &courierRepositoryPostgres{pool: conn}, nil
}

func parseConfig(cfg config.PostgresStorage) (*pgxpool.Config, error) {
	const op = "repository.postgres.parseConfig"

	storageDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DbPort,
		cfg.DbName,
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
