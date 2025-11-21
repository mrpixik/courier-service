package service

import (
	"context"
	"fmt"
	"service-order-avito/internal/adapters"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
)

type CourierRepository interface {
	Create(context.Context, domain.Courier) (int, error)
	GetOneById(context.Context, int) (*domain.Courier, error)
	GetAll(context.Context) ([]domain.Courier, error)
	Update(context.Context, *domain.Courier) error
	DeleteById(context.Context, int) error
}

type courierService struct {
	repository CourierRepository
}

func NewCourierService(repository CourierRepository) *courierService {
	return &courierService{repository: repository}
}

func (cs *courierService) CreateCourier(ctx context.Context, req *dto.CreateCourierRequest) (*dto.CreateCourierResponse, error) {
	if !IsValidName(req.Name) {
		return nil, fmt.Errorf("%w: invalid name", service.ErrInvalidName)
	}
	if !IsValidPhone(req.Phone) {
		return nil, fmt.Errorf("%w: invalid phone", service.ErrInvalidPhone)
	}
	if !IsValidStatus(req.Status) {
		return nil, fmt.Errorf("%w: invalid status", service.ErrInvalidStatus)
	}

	courierDB := domain.Courier{
		Name:   req.Name,
		Phone:  req.Phone,
		Status: req.Status,
	}

	id, err := cs.repository.Create(ctx, courierDB)
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	return &dto.CreateCourierResponse{Id: id}, nil
}

func (cs *courierService) GetCourier(ctx context.Context, req *dto.GetCourierRequest) (*dto.GetCourierResponse, error) {
	courierDb, err := cs.repository.GetOneById(ctx, req.Id)
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	courier := dto.GetCourierResponse{
		Id:     req.Id,
		Name:   courierDb.Name,
		Phone:  courierDb.Phone,
		Status: courierDb.Status,
	}
	return &courier, nil
}

func (cs *courierService) GetAllCouriers(ctx context.Context) ([]dto.GetCourierResponse, error) {
	couriersDb, err := cs.repository.GetAll(ctx)
	if err != nil {

		//return nil, fmt.Errorf(op+": %w: %w", repository.ErrInternalError, err)
		return nil, adapters.ErrUnwrapRepoToService(err)
	}

	couriers := make([]dto.GetCourierResponse, len(couriersDb))
	for i, courierDb := range couriersDb {
		couriers[i] = dto.GetCourierResponse{
			Id:     courierDb.Id,
			Name:   courierDb.Name,
			Phone:  courierDb.Phone,
			Status: courierDb.Status,
		}
	}

	return couriers, nil
}

func (cs *courierService) UpdateCourier(ctx context.Context, req *dto.UpdateCourierRequest) error {
	if req.Name != "" && !IsValidName(req.Name) {
		return service.ErrInvalidName
	}
	if req.Phone != "" && !IsValidPhone(req.Phone) {
		return service.ErrInvalidPhone
	}
	if req.Status != "" && !IsValidStatus(req.Status) {
		return service.ErrInvalidStatus
	}

	courierDb := &domain.Courier{
		Id:     req.Id,
		Name:   req.Name,
		Phone:  req.Phone,
		Status: req.Status,
	}

	if err := cs.repository.Update(ctx, courierDb); err != nil {
		return adapters.ErrUnwrapRepoToService(err)
	}
	return nil
}

func (cs *courierService) DeleteCourier(ctx context.Context, req *dto.DeleteCourierRequest) error {
	err := cs.repository.DeleteById(ctx, req.Id)
	if err != nil {
		return adapters.ErrUnwrapRepoToService(err)
	}

	return nil
}
