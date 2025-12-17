package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/model"
	"strings"
	"time"
)

type courierRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewCourierRepositoryPostgres(pool *pgxpool.Pool) *courierRepositoryPostgres {
	return &courierRepositoryPostgres{pool: pool}
}

// Create создает нового курьера в табличке (с полем транспорт).
func (c *courierRepositoryPostgres) Create(ctx context.Context, courier model.Courier) (int, error) {
	sql := `
        INSERT INTO couriers (name, phone, status, transport_type)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var id int
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		err = tx.QueryRow(ctx, sql, courier.Name, courier.Phone, courier.Status, courier.TransportType).Scan(&id)
	} else { // без транзакции
		err = c.pool.QueryRow(ctx, sql, courier.Name, courier.Phone, courier.Status, courier.TransportType).Scan(&id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return -1, repository.ErrCourierExists
			}
			return -1, repository.ErrInternalError
		}
	}
	return id, err
}

func (c *courierRepositoryPostgres) GetById(ctx context.Context, id int) (model.Courier, error) {
	sql := `
        SELECT name, phone, status, transport_type, created_at, updated_at
        FROM couriers
        WHERE id=$1
    `

	var courier model.Courier
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		err = tx.QueryRow(ctx, sql, id).Scan(
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.TransportType,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		)
	} else { // без транзакции
		err = c.pool.QueryRow(ctx, sql, id).Scan(
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.TransportType,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return model.Courier{}, repository.ErrCourierNotFound
		}
		return model.Courier{}, repository.ErrInternalError
	}

	return courier, err
}

func (c *courierRepositoryPostgres) GetAll(ctx context.Context) ([]model.Courier, error) {
	sql := `
        SELECT id, name, phone, status, transport_type, created_at, updated_at
        FROM couriers
        ORDER BY created_at DESC
    `

	var rows pgx.Rows
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		rows, err = tx.Query(ctx, sql)
	} else { // без транзакции
		rows, err = c.pool.Query(ctx, sql)
	}

	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer rows.Close()

	var couriers []model.Courier
	for rows.Next() {
		var courier model.Courier
		err = rows.Scan(
			&courier.Id,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.TransportType,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		)
		if err != nil {

			//return nil, fmt.Errorf(op+": %w: %w", repository.ErrInternalError, err)
			return nil, repository.ErrInternalError
		}
		couriers = append(couriers, courier)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.ErrInternalError
	}

	return couriers, nil
}

// GetAvailable решил возвращать полный объект domain.Courier, чтобы сохранить логику геттеров.
// Мне кажется, не очень понятно было бы возвращать только id и transport_type, тем более структура достаточно легкая
func (c *courierRepositoryPostgres) GetAvailable(ctx context.Context) (model.Courier, error) {
	sql := `
		SELECT id, name, phone, status, transport_type, total_deliveries, created_at, updated_at
		FROM couriers
		WHERE status = 'available'
		ORDER BY total_deliveries
		LIMIT 1;
    `

	var courier model.Courier
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		err = tx.QueryRow(ctx, sql).Scan(
			&courier.Id,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.TransportType,
			&courier.TotalDeliveries,
			&courier.CreatedAt,
			&courier.UpdatedAt)
	} else { // без транзакции
		err = c.pool.QueryRow(ctx, sql).Scan(
			&courier.Id,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.TransportType,
			&courier.TotalDeliveries,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		)
	}

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {

			return model.Courier{}, repository.ErrNoAvailableCouriers
		}
		return model.Courier{}, repository.ErrInternalError
	}

	return courier, err
}

func (c *courierRepositoryPostgres) Update(ctx context.Context, courier model.Courier) error {
	sqlParts := make([]string, 0)
	fields := make([]interface{}, 0)
	fieldIdx := 1

	// к сожалению, это никак нельзя вынести из репозитория, так как мы обновляем только переданные поля,
	// соответственно, проверять какие передали, придется именно здесь
	// в целом, можно было бы и вынести это в отдельную функцию, но пока это применяется только тут, поэтому не сильно критично
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

	if courier.TransportType != "" {
		sqlParts = append(sqlParts, fmt.Sprintf("transport_type = $%d", fieldIdx))
		fields = append(fields, courier.TransportType)
		fieldIdx++
	}

	if courier.TotalDeliveries != 0 {
		sqlParts = append(sqlParts, fmt.Sprintf("total_deliveries = $%d", fieldIdx))
		fields = append(fields, courier.TotalDeliveries)
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

	var cmdTag pgconn.CommandTag
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		cmdTag, err = tx.Exec(ctx, sql, fields...)
	} else { // без транзакции
		cmdTag, err = c.pool.Exec(ctx, sql, fields...)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return repository.ErrCourierExists
			}
			return repository.ErrInternalError
		}
	}
	if cmdTag.RowsAffected() == 0 {
		return repository.ErrCourierNotFound
	}

	return err
}

func (c *courierRepositoryPostgres) UpdateStatusManyById(ctx context.Context, ids ...int) error {
	sql := `UPDATE couriers
		SET status = 'available'
		WHERE id=ANY($1)
	`

	var cmdTag pgconn.CommandTag
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		cmdTag, err = tx.Exec(ctx, sql, ids)
	} else { // без транзакции
		cmdTag, err = c.pool.Exec(ctx, sql, ids)
	}

	if err != nil {
		return repository.ErrInternalError
	}
	if cmdTag.RowsAffected() == 0 {
		return repository.ErrInternalError
	}

	return nil
}

func (c *courierRepositoryPostgres) DeleteById(ctx context.Context, id int) error {
	sql := `
        DELETE FROM couriers
        WHERE id=$1;
    `

	var cmdTag pgconn.CommandTag
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		cmdTag, err = tx.Exec(ctx, sql, id)
	} else { // без транзакции
		cmdTag, err = c.pool.Exec(ctx, sql, id)
	}

	if err != nil {
		return repository.ErrInternalError
	}
	if cmdTag.RowsAffected() == 0 {
		return repository.ErrCourierNotFound
	}

	return err
}
