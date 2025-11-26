package delivery

import (
	"context"
	"service-order-avito/internal/adapters"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/service/dep"
	"time"
)

type deliveryService struct {
	tm          dep.TransactionManager
	delRepo     dep.DeliveryRepository
	courRepo    dep.CourierRepository
	delTimeCalc dep.DeliveryTimeCalculator
}

func NewDeliveryService(tm dep.TransactionManager, courRepo dep.CourierRepository, delRepo dep.DeliveryRepository) *deliveryService {
	return &deliveryService{tm: tm, delRepo: delRepo, courRepo: courRepo, delTimeCalc: NewDeliveryTimeFactory()}
}

func (ds *deliveryService) AssignDelivery(ctx context.Context, req *dto.AssignDeliveryRequest) (*dto.AssignDeliveryResponse, error) {
	var res *dto.AssignDeliveryResponse
	err := ds.tm.Begin(ctx, func(ctx context.Context) error {
		courier, err := ds.courRepo.GetAvailable(ctx)
		if err != nil {
			return err
		}
		delivery := domain.Delivery{
			CourierId:  courier.Id,
			OrderId:    req.OrderId,
			AssignedAt: time.Now(),
			Deadline:   ds.delTimeCalc.Calculate(courier.TransportType),
		}

		_, err = ds.delRepo.Create(ctx, delivery)
		if err != nil {
			return err
		}

		assignedCourier := domain.Courier{
			Id:              courier.Id,
			Status:          domain.StatusBusy,
			TotalDeliveries: courier.TotalDeliveries + 1,
		}

		err = ds.courRepo.Update(ctx, assignedCourier)
		if err != nil {
			return service.ErrInternalError
		}

		res = &dto.AssignDeliveryResponse{
			CourierId:        courier.Id,
			OrderId:          req.OrderId,
			TransportType:    courier.TransportType,
			DeliveryDeadline: delivery.Deadline,
		}
		return nil
	})
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	return res, nil
}

func (ds *deliveryService) UnassignDelivery(ctx context.Context, req *dto.UnassignDeliveryRequest) (*dto.UnassignDeliveryResponse, error) {
	var res *dto.UnassignDeliveryResponse
	err := ds.tm.Begin(ctx, func(ctx context.Context) error {
		delivery, err := ds.delRepo.GetByOrderId(ctx, req.OrderId)
		if err != nil {
			return err
		}

		err = ds.delRepo.DeleteByOrderId(ctx, req.OrderId)
		if err != nil {
			return err
		}

		courier := domain.Courier{
			Id:     delivery.CourierId,
			Status: domain.StatusAvailable,
		}

		err = ds.courRepo.Update(ctx, courier)
		if err != nil {
			return err
		}

		res = &dto.UnassignDeliveryResponse{
			OrderId:   req.OrderId,
			Status:    "unassigned", // пусть будет эта магическая строчка, а то я не знаю где ее нормально объявить:)
			CourierId: courier.Id,
		}
		return nil
	})
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	return res, nil
}

// UnassignAllCompletedDeliveries завершает все заказы, дедлайн которых прошел, меняет статус ответственных курьеров на 'available'
// возвращает количество завершенных заказов.
func (ds *deliveryService) UnassignAllCompletedDeliveries(ctx context.Context) (int, error) {
	var totalUnassigned int
	err := ds.tm.Begin(ctx, func(ctx context.Context) error {
		completedDeliveries, err := ds.delRepo.GetAllCompleted(ctx)
		if err != nil {
			return err
		}
		if len(completedDeliveries) == 0 {
			return nil
		}

		totalUnassigned = len(completedDeliveries)

		courierIds := make([]int, len(completedDeliveries))
		deliveriesIds := make([]int, len(completedDeliveries))
		for i, d := range completedDeliveries {
			courierIds[i] = d.CourierId
			deliveriesIds[i] = d.Id
		}

		err = ds.delRepo.DeleteManyById(ctx, deliveriesIds...)
		if err != nil {
			return err
		}

		err = ds.courRepo.UpdateStatusManyById(ctx, courierIds...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, adapters.ErrUnwrapRepoToService(err)
	}
	return totalUnassigned, nil
}
