package courier

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/model"
	"service-order-avito/internal/repository/postgres"
	"testing"
	"time"
)

type CourierRepository interface {
	Create(context.Context, model.Courier) (int, error)
	GetById(context.Context, int) (model.Courier, error)
	GetAll(context.Context) ([]model.Courier, error)
	Update(context.Context, model.Courier) error
	UpdateStatusManyById(context.Context, ...int) error
	DeleteById(context.Context, int) error
	GetAvailable(context.Context) (model.Courier, error)
}

type CourierRepositoryTestSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	repo CourierRepository
	ctx  context.Context
}

// TODO: переписать с SetupTest() или TearDownTest()
func TestCourierRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CourierRepositoryTestSuite))
}

func (s *CourierRepositoryTestSuite) SetupSuite() {
	dsn := "postgres://tester:test123@localhost:5433/courier-test-db?sslmode=disable"
	pgxCfg, err := pgxpool.ParseConfig(dsn)
	s.Require().NoError(err)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	s.Require().NoError(err)

	s.pool = pool
	s.repo = postgres.NewCourierRepositoryPostgres(pool)
	s.ctx = context.Background()

	_, err = s.pool.Exec(s.ctx, "DELETE FROM delivery")
	s.Require().NoError(err)

	_, err = s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)
}

func (s *CourierRepositoryTestSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *CourierRepositoryTestSuite) TestCreateCourier_Success() {
	c := model.Courier{
		Name:          "Иван Иванов",
		Phone:         "+79991112233",
		Status:        "available",
		TransportType: "on_foot",
	}

	id, err := s.repo.Create(s.ctx, c)

	s.Require().NoError(err)
	s.Require().Greater(id, 0)

	var count int
	err = s.pool.QueryRow(s.ctx, `SELECT COUNT(*) FROM couriers WHERE id=$1 AND phone=$2`, id, c.Phone).
		Scan(&count)

	s.Require().NoError(err)
	s.Require().Equal(1, count)
}

func (s *CourierRepositoryTestSuite) TestCreateCourier_DuplicatePhone() {
	c := model.Courier{
		Name:          "Пётр Петров",
		Phone:         "+78885556677",
		Status:        "available",
		TransportType: "car",
	}

	id1, err := s.repo.Create(s.ctx, c)
	s.Require().NoError(err)
	s.Require().Greater(id1, 0)

	_, err = s.repo.Create(s.ctx, c)
	s.Require().Error(err)
	s.Require().Equal(repository.ErrCourierExists, err)
}

func (s *CourierRepositoryTestSuite) TestGetById_Success() {
	c := model.Courier{
		Name:          "Тест Тестович",
		Phone:         "+70000000001",
		Status:        "busy",
		TransportType: "car",
	}

	id, err := s.repo.Create(s.ctx, c)
	s.Require().NoError(err)
	s.Require().Greater(id, 0)

	got, err := s.repo.GetById(s.ctx, id)
	s.Require().NoError(err)

	s.Require().Equal(c.Name, got.Name)
	s.Require().Equal(c.Phone, got.Phone)
	s.Require().Equal(c.Status, got.Status)
	s.Require().Equal(c.TransportType, got.TransportType)

	s.Require().False(got.CreatedAt.IsZero())
	s.Require().False(got.UpdatedAt.IsZero())
}

func (s *CourierRepositoryTestSuite) TestGetById_NotFound() {
	_, err := s.repo.GetById(s.ctx, 9999999)
	s.Require().Error(err)
	s.Require().Equal(repository.ErrCourierNotFound, err)
}

func (s *CourierRepositoryTestSuite) TestGetAll_Success() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	c1 := model.Courier{
		Name:          "A1",
		Phone:         "+71111111111",
		Status:        "available",
		TransportType: "on_foot",
	}
	c2 := model.Courier{
		Name:          "A2",
		Phone:         "+72222222222",
		Status:        "busy",
		TransportType: "car",
	}

	_, err = s.repo.Create(s.ctx, c1)
	s.Require().NoError(err)
	time.Sleep(10 * time.Millisecond)

	_, err = s.repo.Create(s.ctx, c2)
	s.Require().NoError(err)

	list, err := s.repo.GetAll(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(list, 2)

	s.Require().Equal("A2", list[0].Name)
	s.Require().Equal("A1", list[1].Name)
}

func (s *CourierRepositoryTestSuite) TestGetAll_Empty() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	list, err := s.repo.GetAll(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(list, 0)
}

func (s *CourierRepositoryTestSuite) TestGetAvailable_Success() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	// Недоступный курьер
	c1 := model.Courier{
		Name:            "Busy",
		Phone:           "+70000000002",
		Status:          "busy",
		TransportType:   "car",
		TotalDeliveries: 10,
	}

	c2 := model.Courier{
		Name:            "Available1",
		Phone:           "+70000000003",
		Status:          "available",
		TransportType:   "on_foot",
		TotalDeliveries: 5,
	}
	c3 := model.Courier{
		Name:            "Available2",
		Phone:           "+70000000004",
		Status:          "available",
		TransportType:   "scooter",
		TotalDeliveries: 1,
	}

	_, err = s.repo.Create(s.ctx, c1)
	s.Require().NoError(err)

	id2, err := s.repo.Create(s.ctx, c2)
	s.Require().NoError(err)

	id3, err := s.repo.Create(s.ctx, c3)
	s.Require().NoError(err)

	_, err = s.pool.Exec(s.ctx, `UPDATE couriers SET total_deliveries=5 WHERE id=$1`, id2)
	s.Require().NoError(err)
	_, err = s.pool.Exec(s.ctx, `UPDATE couriers SET total_deliveries=1 WHERE id=$1`, id3)
	s.Require().NoError(err)

	got, err := s.repo.GetAvailable(s.ctx)
	s.Require().NoError(err)

	s.Require().Equal("Available2", got.Name)
	s.Require().Equal("available", got.Status)
	s.Require().Equal("scooter", got.TransportType)
	s.Require().Equal(1, got.TotalDeliveries)
}

func (s *CourierRepositoryTestSuite) TestGetAvailable_NotFound() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	c := model.Courier{
		Name:          "BusyBusy",
		Phone:         "+78888888888",
		Status:        "busy",
		TransportType: "car",
	}

	_, err = s.repo.Create(s.ctx, c)
	s.Require().NoError(err)

	_, err = s.repo.GetAvailable(s.ctx)
	s.Require().Error(err)
	s.Require().Equal(repository.ErrNoAvailableCouriers, err)
}

func (s *CourierRepositoryTestSuite) TestUpdate_Success() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	c := model.Courier{
		Name:          "Old",
		Phone:         "+75555555555",
		Status:        "busy",
		TransportType: "car",
	}
	id, err := s.repo.Create(s.ctx, c)
	s.Require().NoError(err)

	update := model.Courier{
		Id:              id,
		Name:            "NewName",
		Status:          "available",
		TotalDeliveries: 15,
	}

	err = s.repo.Update(s.ctx, update)
	s.Require().NoError(err)

	var got model.Courier
	err = s.pool.QueryRow(s.ctx,
		`SELECT name, status, total_deliveries FROM couriers WHERE id=$1`, id,
	).Scan(&got.Name, &got.Status, &got.TotalDeliveries)
	s.Require().NoError(err)

	s.Require().Equal("NewName", got.Name)
	s.Require().Equal("available", got.Status)
	s.Require().Equal(15, got.TotalDeliveries)
}

func (s *CourierRepositoryTestSuite) TestUpdate_NotFound() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	update := model.Courier{
		Id:     999999,
		Name:   "New",
		Status: "available",
	}

	err = s.repo.Update(s.ctx, update)
	s.Require().Error(err)
	s.Require().Equal(repository.ErrCourierNotFound, err)
}

func (s *CourierRepositoryTestSuite) TestUpdate_DuplicatePhone() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	// Создаём двух курьеров
	id1, err := s.repo.Create(s.ctx, model.Courier{
		Name:          "A",
		Phone:         "+71111111111",
		Status:        "available",
		TransportType: "car",
	})
	s.Require().NoError(err)

	_, err = s.repo.Create(s.ctx, model.Courier{
		Name:          "B",
		Phone:         "+72222222222",
		Status:        "busy",
		TransportType: "on_foot",
	})
	s.Require().NoError(err)

	err = s.repo.Update(s.ctx, model.Courier{
		Id:    id1,
		Phone: "+72222222222",
	})

	s.Require().Error(err)
	s.Require().Equal(repository.ErrCourierExists, err)
}

func (s *CourierRepositoryTestSuite) TestUpdateStatusManyById_Success() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	id1, _ := s.repo.Create(s.ctx, model.Courier{
		Name: "X1", Phone: "+70000000010", Status: "busy", TransportType: "car",
	})
	id2, _ := s.repo.Create(s.ctx, model.Courier{
		Name: "X2", Phone: "+70000000011", Status: "busy", TransportType: "car",
	})
	_, _ = s.repo.Create(s.ctx, model.Courier{
		Name: "X3", Phone: "+70000000012", Status: "busy", TransportType: "car",
	})

	err = s.repo.UpdateStatusManyById(s.ctx, id1, id2)
	s.Require().NoError(err)

	var statuses []string
	rows, err := s.pool.Query(s.ctx, `SELECT status FROM couriers WHERE id IN ($1, $2)`, id1, id2)
	s.Require().NoError(err)
	defer rows.Close()

	for rows.Next() {
		var st string
		_ = rows.Scan(&st)
		statuses = append(statuses, st)
	}

	s.Require().Contains(statuses, "available")
	s.Require().Len(statuses, 2)
}

func (s *CourierRepositoryTestSuite) TestUpdateStatusManyById_NoRowsAffected() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	err = s.repo.UpdateStatusManyById(s.ctx, 999999, 888888)
	s.Require().Error(err)
	s.Require().Equal(repository.ErrInternalError, err)
}

func (s *CourierRepositoryTestSuite) TestDeleteById_Success() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	id, err := s.repo.Create(s.ctx, model.Courier{
		Name:          "DelTest",
		Phone:         "+73333333333",
		Status:        "available",
		TransportType: "on_foot",
	})
	s.Require().NoError(err)

	err = s.repo.DeleteById(s.ctx, id)
	s.Require().NoError(err)

	var count int
	err = s.pool.QueryRow(s.ctx, `SELECT COUNT(*) FROM couriers WHERE id=$1`, id).
		Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(0, count)
}

func (s *CourierRepositoryTestSuite) TestDeleteById_NotFound() {
	_, err := s.pool.Exec(s.ctx, "DELETE FROM couriers")
	s.Require().NoError(err)

	err = s.repo.DeleteById(s.ctx, 999999)
	s.Require().Error(err)
	s.Require().Equal(repository.ErrCourierNotFound, err)
}
