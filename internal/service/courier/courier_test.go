package courier

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
)

// TODO: можно объединить тесты с ошибками репозитория в 1 табличный
func TestCourierService_CreateCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.CreateCourierRequest{
		Name:          "John",
		Phone:         "+12345678901",
		Status:        domain.StatusAvailable,
		TransportType: "car",
	}

	expectedCourier := domain.Courier{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}

	mockRepo.
		EXPECT().
		Create(gomock.Any(), expectedCourier).
		Return(10, nil)

	resp, err := cs.CreateCourier(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, 10, resp.Id)
}

func TestCourierService_CreateCourier_SuccessUnknownTransportType(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.CreateCourierRequest{
		Name:          "John",
		Phone:         "+12345678901",
		Status:        domain.StatusAvailable,
		TransportType: "balloon",
	}

	expectedCourier := domain.Courier{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: "on_foot",
	}

	mockRepo.
		EXPECT().
		Create(gomock.Any(), expectedCourier).
		Return(10, nil)

	resp, err := cs.CreateCourier(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, 10, resp.Id)
}

func TestCourierService_CreateCourier_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		req         dto.CreateCourierRequest
		expectedErr error
	}{
		{
			name: "invalid name",
			req: dto.CreateCourierRequest{
				Name:          "",
				Phone:         "+12345678901",
				Status:        domain.StatusAvailable,
				TransportType: "car",
			},
			expectedErr: service.ErrInvalidName,
		},
		{
			name: "invalid phone",
			req: dto.CreateCourierRequest{
				Name:          "John",
				Phone:         "bad phone",
				Status:        domain.StatusAvailable,
				TransportType: "car",
			},
			expectedErr: service.ErrInvalidPhone,
		},
		{
			name: "invalid status",
			req: dto.CreateCourierRequest{
				Name:          "John",
				Phone:         "+12345678901",
				Status:        "UNKNOWN",
				TransportType: "car",
			},
			expectedErr: service.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_dep.NewMockCourierRepository(ctrl)
			mockTM := mock_dep.NewMockTransactionManager(ctrl)

			cs := NewCourierService(mockTM, mockRepo)

			resp, err := cs.CreateCourier(context.Background(), &tt.req)

			require.Nil(t, resp)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestCourierService_CreateCourier_CourierExistsError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.CreateCourierRequest{
		Name:          "John",
		Phone:         "+12345678901",
		Status:        domain.StatusAvailable,
		TransportType: "car",
	}

	mockRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(-1, repository.ErrCourierExists)

	_, err := cs.CreateCourier(context.Background(), req)
	require.Equal(t, err, service.ErrCourierExists)
}

func TestCourierService_CreateCourier_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.CreateCourierRequest{
		Name:          "John",
		Phone:         "+12345678901",
		Status:        domain.StatusAvailable,
		TransportType: "car",
	}

	mockRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(-1, repository.ErrInternalError)

	_, err := cs.CreateCourier(context.Background(), req)
	require.Equal(t, err, service.ErrInternalError)
}

func TestCourierService_GetCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.GetCourierRequest{
		Id: 1,
	}

	expectedResponse := &dto.GetCourierResponse{
		Id:            1,
		Name:          "John",
		Phone:         "+12345678901",
		Status:        domain.StatusAvailable,
		TransportType: "car",
	}

	mockRepo.
		EXPECT().
		GetById(gomock.Any(), 1).
		Return(domain.Courier{
			Id:            expectedResponse.Id,
			Name:          expectedResponse.Name,
			Phone:         expectedResponse.Phone,
			Status:        expectedResponse.Status,
			TransportType: expectedResponse.TransportType,
		}, nil)

	resp, err := cs.GetCourier(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, expectedResponse.Id, resp.Id)
	require.Equal(t, expectedResponse.Name, resp.Name)
	require.Equal(t, expectedResponse.Phone, resp.Phone)
	require.Equal(t, expectedResponse.Status, resp.Status)
	require.Equal(t, expectedResponse.TransportType, resp.TransportType)
}

func TestCourierService_GetCourier_CourierNotFoundError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.GetCourierRequest{
		Id: 1,
	}

	mockRepo.
		EXPECT().
		GetById(gomock.Any(), 1).
		Return(domain.Courier{}, repository.ErrCourierNotFound)

	_, err := cs.GetCourier(context.Background(), req)

	require.Equal(t, err, service.ErrCourierNotFound)
}

func TestCourierService_GetCourier_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.GetCourierRequest{
		Id: 1,
	}

	mockRepo.
		EXPECT().
		GetById(gomock.Any(), 1).
		Return(domain.Courier{}, repository.ErrInternalError)

	_, err := cs.GetCourier(context.Background(), req)

	require.Equal(t, err, service.ErrInternalError)
}

func TestCourierService_GetAllCouriers_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	expectedResponse := []dto.GetCourierResponse{
		{
			Id:            1,
			Name:          "John",
			Phone:         "+12345678901",
			Status:        domain.StatusAvailable,
			TransportType: "car",
		},
		{
			Id:            2,
			Name:          "Martin",
			Phone:         "+7777777777",
			Status:        domain.StatusBusy,
			TransportType: "helicopter",
		},
	}

	mockRes := []domain.Courier{
		{
			Id:            1,
			Name:          "John",
			Phone:         "+12345678901",
			Status:        domain.StatusAvailable,
			TransportType: "car",
		},
		{
			Id:            2,
			Name:          "Martin",
			Phone:         "+7777777777",
			Status:        domain.StatusBusy,
			TransportType: "helicopter",
		},
	}

	mockRepo.
		EXPECT().
		GetAll(gomock.Any()).
		Return(mockRes, nil)

	respSlice, err := cs.GetAllCouriers(context.Background())
	require.NoError(t, err)

	for i := range respSlice {
		require.Equal(t, expectedResponse[i].Id, respSlice[i].Id)
		require.Equal(t, expectedResponse[i].Name, respSlice[i].Name)
		require.Equal(t, expectedResponse[i].Phone, respSlice[i].Phone)
		require.Equal(t, expectedResponse[i].Status, respSlice[i].Status)
		require.Equal(t, expectedResponse[i].TransportType, respSlice[i].TransportType)
	}
}

func TestCourierService_GetAllCouriers_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	mockRepo.
		EXPECT().
		GetAll(gomock.Any()).
		Return(nil, repository.ErrInternalError)

	_, err := cs.GetAllCouriers(context.Background())

	require.Equal(t, err, service.ErrInternalError)
}

func TestCourierService_UpdateCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.UpdateCourierRequest{
		Id:            1,
		Name:          "John",
		Phone:         "+12345678901",
		Status:        domain.StatusAvailable,
		TransportType: "car",
	}

	expectedCourier := domain.Courier{
		Id:            req.Id,
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}

	mockRepo.
		EXPECT().
		Update(gomock.Any(), expectedCourier).
		Return(nil)

	err := cs.UpdateCourier(context.Background(), req)
	require.NoError(t, err)
}

func TestCourierService_UpdateCourier_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		req         dto.UpdateCourierRequest
		expectedErr error
	}{
		{
			name: "invalid name",
			req: dto.UpdateCourierRequest{
				Name:          "s1mple",
				Phone:         "+12345678901",
				Status:        domain.StatusAvailable,
				TransportType: "car",
			},
			expectedErr: service.ErrInvalidName,
		},
		{
			name: "invalid phone",
			req: dto.UpdateCourierRequest{
				Name:          "John",
				Phone:         "bad phone",
				Status:        domain.StatusAvailable,
				TransportType: "car",
			},
			expectedErr: service.ErrInvalidPhone,
		},
		{
			name: "invalid status",
			req: dto.UpdateCourierRequest{
				Name:          "John",
				Phone:         "+12345678901",
				Status:        "UNKNOWN",
				TransportType: "car",
			},
			expectedErr: service.ErrInvalidStatus,
		},
		{
			name: "invalid transport type",
			req: dto.UpdateCourierRequest{
				Name:          "John",
				Phone:         "+12345678901",
				Status:        domain.StatusAvailable,
				TransportType: "batmobile",
			},
			expectedErr: service.ErrInvalidTransportType,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_dep.NewMockCourierRepository(ctrl)
			mockTM := mock_dep.NewMockTransactionManager(ctrl)

			cs := NewCourierService(mockTM, mockRepo)

			err := cs.UpdateCourier(context.Background(), &tt.req)

			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestCourierService_UpdateCourier_RepositoryErrors(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.UpdateCourierRequest
		repoErr     error
		expectedErr error
	}{
		{
			name: "courier not found",
			req: &dto.UpdateCourierRequest{
				Id:            1,
				Name:          "John",
				Phone:         "+12345678901",
				Status:        domain.StatusAvailable,
				TransportType: "car",
			},
			repoErr:     repository.ErrCourierNotFound,
			expectedErr: service.ErrCourierNotFound,
		},
		{
			name: "internal error",
			req: &dto.UpdateCourierRequest{
				Id:            1,
				Name:          "John",
				Phone:         "+12345678901",
				Status:        domain.StatusAvailable,
				TransportType: "car",
			},
			repoErr:     repository.ErrInternalError,
			expectedErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_dep.NewMockCourierRepository(ctrl)
			mockTM := mock_dep.NewMockTransactionManager(ctrl)

			cs := NewCourierService(mockTM, mockRepo)

			expectedCourier := domain.Courier{
				Id:            tt.req.Id,
				Name:          tt.req.Name,
				Phone:         tt.req.Phone,
				Status:        tt.req.Status,
				TransportType: tt.req.TransportType,
			}

			mockRepo.
				EXPECT().
				Update(gomock.Any(), expectedCourier).
				Return(tt.repoErr)

			err := cs.UpdateCourier(context.Background(), tt.req)
			require.Equal(t, err, tt.expectedErr)
		})
	}

}

func TestCourierService_DeleteCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_dep.NewMockCourierRepository(ctrl)
	mockTM := mock_dep.NewMockTransactionManager(ctrl)

	cs := NewCourierService(mockTM, mockRepo)

	req := &dto.DeleteCourierRequest{
		Id: 1,
	}

	mockRepo.
		EXPECT().
		DeleteById(gomock.Any(), 1).
		Return(nil)

	err := cs.DeleteCourier(context.Background(), req)
	require.NoError(t, err)
}

func TestCourierService_DeleteCourier_RepositoryErrors(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.DeleteCourierRequest
		repoErr     error
		expectedErr error
	}{
		{
			name: "courier not found",
			req: &dto.DeleteCourierRequest{
				Id: 1,
			},
			repoErr:     repository.ErrCourierNotFound,
			expectedErr: service.ErrCourierNotFound,
		},
		{
			name: "internal error",
			req: &dto.DeleteCourierRequest{
				Id: 1,
			},
			repoErr:     repository.ErrInternalError,
			expectedErr: service.ErrInternalError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_dep.NewMockCourierRepository(ctrl)
			mockTM := mock_dep.NewMockTransactionManager(ctrl)

			cs := NewCourierService(mockTM, mockRepo)

			mockRepo.
				EXPECT().
				DeleteById(gomock.Any(), tt.req.Id).
				Return(tt.repoErr)

			err := cs.DeleteCourier(context.Background(), tt.req)
			require.Equal(t, err, tt.expectedErr)
		})
	}

}
