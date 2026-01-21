package delivery

import (
	"context"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/model"
	"service-order-avito/internal/repository/postgres"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type DeliveryRepository interface {
	Create(context.Context, model.Delivery) (int, error)
	GetByOrderId(context.Context, string) (model.Delivery, error)
	GetAllCompleted(context.Context) ([]model.Delivery, error)
	DeleteByOrderId(context.Context, string) error
	DeleteManyById(context.Context, ...int) error
}

type DeliveryRepositoryTestSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	repo DeliveryRepository
	ctx  context.Context
}

func TestDeliveryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryRepositoryTestSuite))
}

func (s *DeliveryRepositoryTestSuite) SetupSuite() {
	dsn := "postgres://tester:test123@localhost:5433/courier-test-db?sslmode=disable"
	pgxCfg, err := pgxpool.ParseConfig(dsn)
	s.Require().NoError(err)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	s.Require().NoError(err)

	s.pool = pool
	s.repo = postgres.NewDeliveryRepositoryPostgres(pool)
	s.ctx = context.Background()
}

func (s *DeliveryRepositoryTestSuite) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE delivery, couriers RESTART IDENTITY CASCADE")
	s.Require().NoError(err)
}

func (s *DeliveryRepositoryTestSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *DeliveryRepositoryTestSuite) TestCreate_Success() {
	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO couriers (id, name, phone, status, transport_type, total_deliveries, created_at)
        VALUES (1, 'John', '123', 'available', 'car', 0, NOW())
    `)
	s.Require().NoError(err)

	delivery := model.Delivery{
		CourierId:  1,
		OrderId:    "order-1",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(2 * time.Hour),
	}

	id, err := s.repo.Create(s.ctx, delivery)

	s.Require().NoError(err)
	s.Assert().Greater(id, 0)

	var count int
	err = s.pool.QueryRow(s.ctx,
		`SELECT COUNT(*) FROM delivery WHERE id=$1 AND order_id=$2`,
		id, "order-1",
	).Scan(&count)

	s.Require().NoError(err)
	s.Assert().Equal(1, count)
}

func (s *DeliveryRepositoryTestSuite) TestCreate_DuplicateOrderId() {
	// создаём курьера
	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO couriers (id, name, phone, status, transport_type, total_deliveries, created_at)
        VALUES (2, 'Mike', '555', 'available', 'bike', 0, NOW())
    `)
	s.Require().NoError(err)

	delivery := model.Delivery{
		CourierId:  2,
		OrderId:    "duplicate-order",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(1 * time.Hour),
	}

	_, err = s.repo.Create(s.ctx, delivery)
	s.Require().NoError(err)

	_, err = s.repo.Create(s.ctx, delivery)

	s.Require().Error(err)
	s.Assert().ErrorIs(err, repository.ErrDeliveryExists)
}

func (s *DeliveryRepositoryTestSuite) TestCreate_InvalidCourierId() {
	delivery := model.Delivery{
		CourierId:  9999,
		OrderId:    "invalid-fk-order",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(2 * time.Hour),
	}

	id, err := s.repo.Create(s.ctx, delivery)

	s.Require().Error(err)
	s.Assert().Equal(-1, id)
	s.Assert().ErrorIs(err, repository.ErrInternalError)
}

func (s *DeliveryRepositoryTestSuite) TestGetByOrderId_Success() {
	// создаём курьера
	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO couriers (id, name, phone, status, transport_type, total_deliveries, created_at)
        VALUES (1, 'Mike', '555', 'available', 'bike', 0, NOW())
    `)
	s.Require().NoError(err)

	assigned := time.Now().Add(-10 * time.Minute)
	deadline := time.Now().Add(10 * time.Minute)

	id := s.insertDelivery(1, "ORDER-XYZ", assigned, deadline)
	s.Require().Greater(id, 0)

	d, err := s.repo.GetByOrderId(s.ctx, "ORDER-XYZ")
	s.Require().NoError(err)
	s.Equal(id, d.Id)
	s.Equal("ORDER-XYZ", d.OrderId)
}

func (s *DeliveryRepositoryTestSuite) TestGetByOrderId_NotFound() {
	_, err := s.repo.GetByOrderId(s.ctx, "NOT_EXIST")
	s.Require().ErrorIs(err, repository.ErrDeliveryNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestGetAllCompleted_ReturnsPastDeadlines() {
	// создаём курьера
	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO couriers (id, name, phone, status, transport_type, total_deliveries, created_at)
        VALUES (1, 'Mike', '555', 'available', 'bike', 0, NOW())
    `)
	s.Require().NoError(err)

	s.insertDelivery(1, "O1", time.Now().Add(-1*time.Hour), time.Now().Add(-30*time.Minute))
	s.insertDelivery(1, "O2", time.Now().Add(-2*time.Hour), time.Now().Add(-10*time.Minute))

	s.insertDelivery(1, "O3", time.Now(), time.Now().Add(1*time.Hour))

	list, err := s.repo.GetAllCompleted(s.ctx)
	s.Require().NoError(err)

	s.Len(list, 2)
	s.ElementsMatch(
		[]string{"O1", "O2"},
		[]string{list[0].OrderId, list[1].OrderId},
	)
}

func (s *DeliveryRepositoryTestSuite) TestDeleteByOrderId_Success() {
	// создаём курьера
	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO couriers (id, name, phone, status, transport_type, total_deliveries, created_at)
        VALUES (1, 'Mike', '555', 'available', 'bike', 0, NOW())
    `)
	s.Require().NoError(err)

	delivery := model.Delivery{
		CourierId:  1,
		OrderId:    "DEL-99",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(2 * time.Hour),
	}

	_, err = s.repo.Create(s.ctx, delivery)

	s.Require().NoError(err)

	err = s.repo.DeleteByOrderId(s.ctx, "DEL-99")
	s.Require().NoError(err)

	_, err = s.repo.GetByOrderId(s.ctx, "DEL-99")
	s.ErrorIs(err, repository.ErrDeliveryNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestDeleteByOrderId_NotFound() {
	err := s.repo.DeleteByOrderId(s.ctx, "NOPE")
	s.ErrorIs(err, repository.ErrDeliveryNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestDeleteManyById_Success() {
	// создаём курьера
	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO couriers (id, name, phone, status, transport_type, total_deliveries, created_at)
        VALUES (1, 'Mike', '555', 'available', 'bike', 0, NOW())
    `)
	s.Require().NoError(err)

	id1 := s.insertDelivery(1, "B1", time.Now(), time.Now())
	id2 := s.insertDelivery(1, "B2", time.Now(), time.Now())

	err = s.repo.DeleteManyById(s.ctx, id1, id2)
	s.Require().NoError(err)

	_, err = s.repo.GetByOrderId(s.ctx, "B1")
	s.ErrorIs(err, repository.ErrDeliveryNotFound)

	_, err = s.repo.GetByOrderId(s.ctx, "B2")
	s.ErrorIs(err, repository.ErrDeliveryNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestDeleteManyById_NotFound() {
	err := s.repo.DeleteManyById(s.ctx, 999999)
	s.ErrorIs(err, repository.ErrDeliveryNotFound)
}

func (s *DeliveryRepositoryTestSuite) insertDelivery(courierId int, orderId string, assigned, deadline time.Time) int {
	var id int
	err := s.pool.QueryRow(s.ctx,
		`INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		courierId, orderId, assigned, deadline,
	).Scan(&id)

	s.Require().NoError(err)
	return id
}
