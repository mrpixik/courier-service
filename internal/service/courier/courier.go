package courier

import (
	"context"
	"service-order-avito/internal/adapters"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/service/dep"
)

type courierService struct {
	tm         dep.TransactionManager
	repository dep.CourierRepository
}

func NewCourierService(tm dep.TransactionManager, repository dep.CourierRepository) *courierService {
	return &courierService{tm: tm, repository: repository}
}

func (cs *courierService) CreateCourier(ctx context.Context, req *dto.CreateCourierRequest) (*dto.CreateCourierResponse, error) {
	if !IsValidName(req.Name) {
		return nil, service.ErrInvalidName
	}
	if !IsValidPhone(req.Phone) {
		return nil, service.ErrInvalidPhone
	}
	if !IsValidStatus(req.Status) {
		return nil, service.ErrInvalidStatus
	}
	// выбрал вариант не возвращать ошибку, так как все равно есть дефолтное значение "on_foot"
	// (хотя по такой логике, надо было и с полем status так же сделать)
	if !IsValidTransportType(req.TransportType) {
		req.TransportType = "on_foot"
	}

	courierDB := domain.Courier{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}

	id, err := cs.repository.Create(ctx, courierDB)

	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	return &dto.CreateCourierResponse{Id: id}, nil
}

func (cs *courierService) GetCourier(ctx context.Context, req *dto.GetCourierRequest) (*dto.GetCourierResponse, error) {
	courierDb, err := cs.repository.GetById(ctx, req.Id)
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}

	courier := dto.GetCourierResponse{
		Id:            req.Id,
		Name:          courierDb.Name,
		Phone:         courierDb.Phone,
		Status:        courierDb.Status,
		TransportType: courierDb.TransportType,
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
			Id:            courierDb.Id,
			Name:          courierDb.Name,
			Phone:         courierDb.Phone,
			Status:        courierDb.Status,
			TransportType: courierDb.TransportType,
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
	if req.TransportType != "" && !IsValidTransportType(req.TransportType) {
		return service.ErrInvalidTransportType
	}

	courierDb := domain.Courier{
		Id:            req.Id,
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
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
