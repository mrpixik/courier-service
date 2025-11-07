package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/http/server/dto"
)

type CourierRepository interface {
	Create(context.Context, domain.Courier) (int, error)
	GetOneById(context.Context, int) (*domain.Courier, error)
	GetAll(context.Context) ([]domain.Courier, error)
	Update(context.Context, *domain.Courier) error
}

// Хочу добавить сюда логгер для того, чтобы отслеживать неявное поведение бд
// (в случае если она вдруг падает или проападет соединение с ней)
type courierService struct {
	repository CourierRepository
}

func NewCourierService(repository CourierRepository) *courierService {
	return &courierService{repository: repository}
}

func (cs *courierService) CreateCourier(ctx context.Context, courierReq *dto.CourierCreateRequest) (int, error) {
	if !IsValidName(courierReq.Name) {
		return -1, fmt.Errorf("%w: invalid name", ErrInvalidName)
	}
	if !IsValidPhone(courierReq.Phone) {
		return -1, fmt.Errorf("%w: invalid phone", ErrInvalidPhone)
	}
	if !IsValidStatus(courierReq.Status) {
		return -1, fmt.Errorf("%w: invalid status", ErrInvalidStatus)
	}

	courierDB := domain.Courier{
		Name:   courierReq.Name,
		Phone:  courierReq.Phone,
		Status: courierReq.Status,
	}

	id, err := cs.repository.Create(ctx, courierDB)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return -1, ErrCourierAlreadyExists
			}
			// вот такого рода непредвиденные ошибки было бы удобно логгировать
			return -1, ErrInternalError
		}
	}
	return id, nil
}

func (cs *courierService) GetCourier(ctx context.Context, id int) (*dto.Courier, error) {
	courierDb, err := cs.repository.GetOneById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCourierNotFound
		}
		return nil, ErrInternalError
	}
	courier := dto.Courier{
		Id:     id,
		Name:   courierDb.Name,
		Phone:  courierDb.Phone,
		Status: courierDb.Status,
	}
	return &courier, nil
}

func (cs *courierService) GetAllCouriers(ctx context.Context) ([]dto.Courier, error) {
	couriersDb, err := cs.repository.GetAll(ctx)
	if err != nil {

		//return nil, fmt.Errorf(op+": %w: %w", repository.ErrUnknownError, err)
		return nil, ErrInternalError
	}
	couriers := make([]dto.Courier, len(couriersDb))
	for i, courierDb := range couriersDb {
		couriers[i] = dto.Courier{
			Id:     courierDb.Id,
			Name:   courierDb.Name,
			Phone:  courierDb.Phone,
			Status: courierDb.Status,
		}
	}

	return couriers, nil
}

func (cs *courierService) UpdateCourier(ctx context.Context, courierReq *dto.CourierUpdateRequest) error {
	if courierReq.Name != "" && !IsValidName(courierReq.Name) {
		return ErrInvalidName
	}
	if courierReq.Phone != "" && !IsValidPhone(courierReq.Phone) {
		return ErrInvalidPhone
	}
	if courierReq.Status != "" && !IsValidStatus(courierReq.Status) {
		return ErrInvalidStatus
	}

	courierDb := &domain.Courier{
		Id:     courierReq.Id,
		Name:   courierReq.Name,
		Phone:  courierReq.Phone,
		Status: courierReq.Status,
	}

	if err := cs.repository.Update(ctx, courierDb); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCourierNotFound
		}
		return err
	}
	return nil
}
