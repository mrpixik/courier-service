package postgres

import (
	"context"
	"fmt"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/repository"
	"strings"
	"time"
)

func (s *courierRepositoryPostgres) CloseConnection() {
	s.pool.Close()
	return
}

func (s *courierRepositoryPostgres) Create(ctx context.Context, courier domain.Courier) (int, error) {
	sql := `
        INSERT INTO couriers (name, phone, status)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	var id int
	err := s.pool.QueryRow(ctx, sql, courier.Name, courier.Phone, courier.Status).Scan(&id)

	return id, err
}

func (s *courierRepositoryPostgres) GetOneById(ctx context.Context, id int) (*domain.Courier, error) {
	sql := `
        SELECT name, phone, status, created_at, updated_at
        FROM couriers
        WHERE id=$1
    `
	var courier domain.Courier
	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.CreatedAt,
		&courier.UpdatedAt,
	)

	return &courier, err
}

func (s *courierRepositoryPostgres) GetAll(ctx context.Context) ([]domain.Courier, error) {
	sql := `
        SELECT id, name, phone, status, created_at, updated_at
        FROM couriers
        ORDER BY created_at DESC
    `

	rows, err := s.pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %w", err)
	}
	defer rows.Close()

	var couriers []domain.Courier
	for rows.Next() {
		var courier domain.Courier
		err := rows.Scan(
			&courier.Id,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", repository.ErrUnknownError, err)
		}
		couriers = append(couriers, courier)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating couriers rows: %w", err)
	}

	return couriers, nil
}

func (s *courierRepositoryPostgres) Update(ctx context.Context, courier *domain.Courier) error {
	sqlParts := make([]string, 0)
	fields := make([]interface{}, 0)
	fieldIdx := 1

	if courier.Name != "" {
		sqlParts = append(sqlParts, fmt.Sprintf("name = $%d", fieldIdx))
		fields = append(fields, courier.Name)
		fieldIdx++
	}

	if courier.Phone != "" {
		sqlParts = append(sqlParts, fmt.Sprintf("phone = $%d", fieldIdx))
		fields = append(fields, courier.Phone)
		fieldIdx++
	}

	if courier.Status != "" {
		sqlParts = append(sqlParts, fmt.Sprintf("status = $%d", fieldIdx))
		fields = append(fields, courier.Status)
		fieldIdx++
	}

	if len(sqlParts) == 0 {
		return nil
	}

	sqlParts = append(sqlParts, fmt.Sprintf("updated_at = $%d", fieldIdx))
	fields = append(fields, time.Now())
	fieldIdx++

	sql := `UPDATE couriers SET` + " " + strings.Join(sqlParts, ", ") + fmt.Sprintf(` WHERE id = $%d`, fieldIdx)
	fields = append(fields, courier.Id)

	_, err := s.pool.Exec(ctx, sql, fields...)

	return err
}
