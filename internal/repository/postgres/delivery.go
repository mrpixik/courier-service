package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/errors/repository"
	"time"
)

type deliveryRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepositoryPostgres(pool *pgxpool.Pool) *deliveryRepositoryPostgres {
	return &deliveryRepositoryPostgres{pool: pool}
}

func (d *deliveryRepositoryPostgres) Create(ctx context.Context, delivery domain.Delivery) (int, error) {
	sql := `
        INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id int
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		err = tx.QueryRow(ctx, sql, delivery.CourierId, delivery.OrderId, delivery.AssignedAt, delivery.Deadline).Scan(&id)
	} else { // без транзакции
		err = d.pool.QueryRow(ctx, sql, delivery.CourierId, delivery.OrderId, delivery.AssignedAt, delivery.Deadline).Scan(&id)
	}

	if err != nil {
		fmt.Println(err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return -1, repository.ErrDeliveryExists
			}
			return -1, repository.ErrInternalError
		}
	}
	return id, err
}

func (c *deliveryRepositoryPostgres) GetByOrderId(ctx context.Context, orderId string) (domain.Delivery, error) {
	sql := `
        SELECT id, courier_id, order_id, assigned_at, deadline
        FROM delivery
        WHERE order_id=$1
    `

	var delivery domain.Delivery
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		err = tx.QueryRow(ctx, sql, orderId).Scan(
			&delivery.Id,
			&delivery.CourierId,
			&delivery.OrderId,
			&delivery.AssignedAt,
			&delivery.Deadline,
		)
	} else { // без транзакции
		err = c.pool.QueryRow(ctx, sql, orderId).Scan(
			&delivery.Id,
			&delivery.CourierId,
			&delivery.OrderId,
			&delivery.AssignedAt,
			&delivery.Deadline,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return domain.Delivery{}, repository.ErrDeliveryNotFound
		}
		return domain.Delivery{}, repository.ErrInternalError
	}

	return delivery, err
}

func (c *deliveryRepositoryPostgres) GetAllCompleted(ctx context.Context) ([]domain.Delivery, error) {
	sql := `
        SELECT id, courier_id, order_id, assigned_at, deadline
        FROM delivery
        WHERE deadline<$1
    `

	var rows pgx.Rows
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		rows, err = tx.Query(ctx, sql, time.Now())
	} else { // без транзакции
		rows, err = c.pool.Query(ctx, sql, time.Now())
	}

	if err != nil {
		return nil, repository.ErrInternalError
	}
	defer rows.Close()

	var deliveries []domain.Delivery
	for rows.Next() {
		var delivery domain.Delivery
		err = rows.Scan(
			&delivery.Id,
			&delivery.CourierId,
			&delivery.OrderId,
			&delivery.AssignedAt,
			&delivery.Deadline,
		)
		if err != nil {
			return nil, repository.ErrInternalError
		}
		deliveries = append(deliveries, delivery)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.ErrInternalError
	}

	return deliveries, nil
}

func (d *deliveryRepositoryPostgres) DeleteByOrderId(ctx context.Context, orderId string) error {
	sql := `
        DELETE FROM delivery
        WHERE order_id=$1
    `

	var cmdTag pgconn.CommandTag
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		cmdTag, err = tx.Exec(ctx, sql, orderId)
	} else { // без транзакции
		cmdTag, err = d.pool.Exec(ctx, sql, orderId)
	}

	if err != nil {
		return repository.ErrInternalError
	}
	if cmdTag.RowsAffected() == 0 {
		return repository.ErrDeliveryNotFound
	}

	return err
}

// решил массовое удаление реализовать через id, а не orderId. так как логично, что этот метод будет использоваться только внутренними обработчиками микросервиса
// в случае, если потребуется добавить функционал для массового удаления доставок по запросу извне, ничего не мешает написать новый метод.
// мне кажется, что мой подход оптимален, так как метод будет вызываться очень часто (ро логике)
// и строка orderId длиной 36 символов, что уже делает разницу существенной, не говоря уже о том, что их будет передаваться сразу несколько
func (d *deliveryRepositoryPostgres) DeleteManyById(ctx context.Context, ids ...int) error {
	sql := `
        DELETE FROM delivery
        WHERE id=ANY($1)
    `

	var cmdTag pgconn.CommandTag
	var err error

	if tx := GetTx(ctx); tx != nil { // с транзакцией
		cmdTag, err = tx.Exec(ctx, sql, ids)
	} else { // без транзакции
		cmdTag, err = d.pool.Exec(ctx, sql, ids)
	}

	if err != nil {
		return repository.ErrInternalError
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrDeliveryNotFound
	}

	return err
}
