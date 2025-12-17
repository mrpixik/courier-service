package delivery

import (
	"context"
	"service-order-avito/internal/adapters"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/domain/model"
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

func (ds *deliveryService) Assign(ctx context.Context, req *dto.AssignDeliveryRequest) (*dto.AssignDeliveryResponse, error) {
	var res *dto.AssignDeliveryResponse
	err := ds.tm.Begin(ctx, func(ctx context.Context) error {
		courier, err := ds.courRepo.GetAvailable(ctx)
		if err != nil {
			return err
		}
		delivery := model.Delivery{
			CourierId:  courier.Id,
			OrderId:    req.OrderId,
			AssignedAt: time.Now(),
			Deadline:   ds.delTimeCalc.Calculate(courier.TransportType),
		}

		_, err = ds.delRepo.Create(ctx, delivery)
		if err != nil {
			return err
		}

		assignedCourier := model.Courier{
			Id:              courier.Id,
			Status:          model.StatusBusy,
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

func (ds *deliveryService) Unassign(ctx context.Context, req *dto.UnassignDeliveryRequest) (*dto.UnassignDeliveryResponse, error) {
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

		courier := model.Courier{
			Id:     delivery.CourierId,
			Status: model.StatusAvailable,
		}

		err = ds.courRepo.Update(ctx, courier)
		if err != nil {
			return err
		}

		res = &dto.UnassignDeliveryResponse{
			OrderId:   req.OrderId,
			Status:    model.StatusUnassigned,
			CourierId: courier.Id,
		}
		return nil
	})
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	return res, nil
}

// UnassignAllCompleted завершает все заказы, дедлайн которых прошел, меняет статус ответственных курьеров на 'available'
// возвращает количество завершенных заказов.
func (ds *deliveryService) UnassignAllCompleted(ctx context.Context) (int, error) {
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

// Complete завершает заказ. Освобождает курьера, но не удаляет заказ из таблицы delivery
func (ds *deliveryService) Complete(ctx context.Context, req *dto.CompleteDeliveryRequest) (*dto.CompleteDeliveryResponse, error) {
	var res *dto.CompleteDeliveryResponse
	err := ds.tm.Begin(ctx, func(ctx context.Context) error {
		delivery, err := ds.delRepo.GetByOrderId(ctx, req.OrderId)
		if err != nil {
			return err
		}

		courier := model.Courier{
			Id:     delivery.CourierId,
			Status: model.StatusAvailable,
		}

		err = ds.courRepo.Update(ctx, courier)
		if err != nil {
			return err
		}

		res = &dto.CompleteDeliveryResponse{
			OrderId:   req.OrderId,
			Status:    model.StatusCompleted,
			CourierId: courier.Id,
		}
		return nil
	})
	if err != nil {
		return nil, adapters.ErrUnwrapRepoToService(err)
	}
	return res, nil
}
