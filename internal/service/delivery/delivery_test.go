package delivery

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/service"
	mock_dep "service-order-avito/internal/service/dep/mocks"
	"testing"
	"time"
)

func TestDeliveryService_AssignDelivery_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)

	ctx := context.Background()
	req := &dto.AssignDeliveryRequest{
		OrderId: "ORDER-123",
	}

	courier := domain.Courier{
		Id:              1,
		Name:            "John",
		Phone:           "+12345678901",
		Status:          domain.StatusAvailable,
		TransportType:   "car",
		TotalDeliveries: 0,
	}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	mockCourierRepo.EXPECT().GetAvailable(gomock.Any()).Return(courier, nil)

	mockDeliveryRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, delivery domain.Delivery) (int, error) {

			require.Equal(t, req.OrderId, delivery.OrderId)
			require.Equal(t, courier.Id, delivery.CourierId)
			require.WithinDuration(t, time.Now(), delivery.AssignedAt, time.Second)
			return 1, nil
		},
	)

	mockCourierRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, updated domain.Courier) error {
			require.Equal(t, courier.Id, updated.Id)
			require.Equal(t, domain.StatusBusy, updated.Status)
			require.Equal(t, courier.TotalDeliveries+1, updated.TotalDeliveries)
			return nil
		},
	)

	resp, err := ds.AssignDelivery(ctx, req)

	require.NoError(t, err)
	require.Equal(t, courier.Id, resp.CourierId)
	require.Equal(t, req.OrderId, resp.OrderId)
	require.Equal(t, courier.TransportType, resp.TransportType)
	require.WithinDuration(t, time.Now(), resp.DeliveryDeadline, 2*time.Hour) // допустимо ±2 часа
}
func TestDeliveryService_AssignDelivery_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)

	ctx := context.Background()
	req := &dto.AssignDeliveryRequest{OrderId: "ORDER-123"}

	t.Run("courier repo GetAvailable returns no available couriers", func(t *testing.T) {
		mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		)
		mockCourierRepo.EXPECT().GetAvailable(gomock.Any()).Return(domain.Courier{}, repository.ErrNoAvailableCouriers)

		resp, err := ds.AssignDelivery(ctx, req)
		require.Nil(t, resp)
		require.ErrorIs(t, err, service.ErrNoAvailableCouriers)
	})

	t.Run("delivery repo Create returns ErrDeliveryExists", func(t *testing.T) {
		courier := domain.Courier{Id: 1, Status: domain.StatusAvailable, TransportType: "car", TotalDeliveries: 0}
		mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		)
		mockCourierRepo.EXPECT().GetAvailable(gomock.Any()).Return(courier, nil)
		mockDeliveryRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(0, repository.ErrDeliveryExists)

		resp, err := ds.AssignDelivery(ctx, req)
		require.Nil(t, resp)
		require.ErrorIs(t, err, service.ErrDeliveryExists)
	})

	t.Run("courier repo Update returns internal error", func(t *testing.T) {
		courier := domain.Courier{Id: 1, Status: domain.StatusAvailable, TransportType: "car", TotalDeliveries: 0}
		mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		)
		mockCourierRepo.EXPECT().GetAvailable(gomock.Any()).Return(courier, nil)
		mockDeliveryRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil)
		mockCourierRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(service.ErrInternalError)

		resp, err := ds.AssignDelivery(ctx, req)
		require.Nil(t, resp)
		require.ErrorIs(t, err, service.ErrInternalError)
	})
}

func TestDeliveryService_UnassignDelivery_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)

	ctx := context.Background()
	req := &dto.UnassignDeliveryRequest{
		OrderId: "ORDER-123",
	}

	delivery := domain.Delivery{
		CourierId: 1,
		OrderId:   req.OrderId,
		Deadline:  time.Now().Add(2 * time.Hour),
	}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		},
	)

	mockDeliveryRepo.EXPECT().GetByOrderId(gomock.Any(), req.OrderId).Return(delivery, nil)
	mockDeliveryRepo.EXPECT().DeleteByOrderId(gomock.Any(), req.OrderId).Return(nil)
	mockCourierRepo.EXPECT().Update(gomock.Any(), domain.Courier{
		Id:     delivery.CourierId,
		Status: domain.StatusAvailable,
	}).Return(nil)

	resp, err := ds.UnassignDelivery(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, delivery.CourierId, resp.CourierId)
	require.Equal(t, req.OrderId, resp.OrderId)
	require.Equal(t, "unassigned", resp.Status)
}

func TestDeliveryService_UnassignDelivery_DeliveryNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()
	req := &dto.UnassignDeliveryRequest{OrderId: "ORDER-123"}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)
	mockDeliveryRepo.EXPECT().GetByOrderId(gomock.Any(), req.OrderId).Return(domain.Delivery{}, repository.ErrDeliveryNotFound)

	resp, err := ds.UnassignDelivery(ctx, req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, service.ErrDeliveryNotFound)

}

func TestDeliveryService_UnassignDelivery_DeleteFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()
	req := &dto.UnassignDeliveryRequest{OrderId: "ORDER-123"}

	delivery := domain.Delivery{CourierId: 1, OrderId: req.OrderId}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)
	mockDeliveryRepo.EXPECT().GetByOrderId(gomock.Any(), req.OrderId).Return(delivery, nil)
	mockDeliveryRepo.EXPECT().DeleteByOrderId(gomock.Any(), req.OrderId).Return(repository.ErrInternalError)

	resp, err := ds.UnassignDelivery(ctx, req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, service.ErrInternalError)

}

func TestDeliveryService_UnassignDelivery_CourierUpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()
	req := &dto.UnassignDeliveryRequest{OrderId: "ORDER-123"}

	delivery := domain.Delivery{CourierId: 1, OrderId: req.OrderId}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)
	mockDeliveryRepo.EXPECT().GetByOrderId(gomock.Any(), req.OrderId).Return(delivery, nil)
	mockDeliveryRepo.EXPECT().DeleteByOrderId(gomock.Any(), req.OrderId).Return(nil)
	mockCourierRepo.EXPECT().Update(gomock.Any(), domain.Courier{
		Id:     delivery.CourierId,
		Status: domain.StatusAvailable,
	}).Return(repository.ErrInternalError)

	resp, err := ds.UnassignDelivery(ctx, req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, service.ErrInternalError)

}

func TestDeliveryService_UnassignAllCompletedDeliveries_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()

	completedDeliveries := []domain.Delivery{
		{Id: 1, CourierId: 101, OrderId: "ORDER-1"},
		{Id: 2, CourierId: 102, OrderId: "ORDER-2"},
	}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)

	mockDeliveryRepo.EXPECT().GetAllCompleted(gomock.Any()).Return(completedDeliveries, nil)
	mockDeliveryRepo.EXPECT().DeleteManyById(gomock.Any(), 1, 2).Return(nil)
	mockCourierRepo.EXPECT().UpdateStatusManyById(gomock.Any(), 101, 102).Return(nil)

	total, err := ds.UnassignAllCompletedDeliveries(ctx)

	require.NoError(t, err)
	require.Equal(t, len(completedDeliveries), total)
}

func TestDeliveryService_UnassignAllCompletedDeliveries_NoCompletedDeliveries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)

	// Возвращаем пустой слайс completedDeliveries
	mockDeliveryRepo.EXPECT().GetAllCompleted(gomock.Any()).Return([]domain.Delivery{}, nil)

	total, err := ds.UnassignAllCompletedDeliveries(ctx)

	require.NoError(t, err)
	require.Equal(t, 0, total)

}

func TestDeliveryService_UnassignAllCompletedDeliveries_GetAllCompletedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)

	mockDeliveryRepo.EXPECT().GetAllCompleted(gomock.Any()).Return(nil, repository.ErrInternalError)

	total, err := ds.UnassignAllCompletedDeliveries(ctx)

	require.Equal(t, 0, total)
	require.ErrorIs(t, err, service.ErrInternalError)

}

func TestDeliveryService_UnassignAllCompletedDeliveries_DeleteManyByIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()

	completedDeliveries := []domain.Delivery{
		{Id: 1, CourierId: 101, OrderId: "ORDER-1"},
		{Id: 2, CourierId: 102, OrderId: "ORDER-2"},
	}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)

	mockDeliveryRepo.EXPECT().GetAllCompleted(gomock.Any()).Return(completedDeliveries, nil)
	mockDeliveryRepo.EXPECT().DeleteManyById(gomock.Any(), 1, 2).Return(repository.ErrDeliveryNotFound)

	total, err := ds.UnassignAllCompletedDeliveries(ctx)

	require.Equal(t, 0, total)
	require.ErrorIs(t, err, service.ErrDeliveryNotFound)

}

func TestDeliveryService_UnassignAllCompletedDeliveries_UpdateStatusManyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTM := mock_dep.NewMockTransactionManager(ctrl)
	mockCourierRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockDeliveryRepo := mock_dep.NewMockDeliveryRepository(ctrl)

	ds := NewDeliveryService(mockTM, mockCourierRepo, mockDeliveryRepo)
	ctx := context.Background()

	completedDeliveries := []domain.Delivery{
		{Id: 1, CourierId: 101, OrderId: "ORDER-1"},
		{Id: 2, CourierId: 102, OrderId: "ORDER-2"},
	}

	mockTM.EXPECT().Begin(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) },
	)

	mockDeliveryRepo.EXPECT().GetAllCompleted(gomock.Any()).Return(completedDeliveries, nil)
	mockDeliveryRepo.EXPECT().DeleteManyById(gomock.Any(), 1, 2).Return(nil)
	mockCourierRepo.EXPECT().UpdateStatusManyById(gomock.Any(), 101, 102).Return(repository.ErrInternalError)

	total, err := ds.UnassignAllCompletedDeliveries(ctx)

	require.Equal(t, 0, total)
	require.ErrorIs(t, err, service.ErrInternalError)
}
