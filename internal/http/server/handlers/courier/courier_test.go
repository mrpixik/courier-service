package courier

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
	"service-order-avito/internal/domain/errors/service"
	mock_courier "service-order-avito/internal/http/server/handlers/courier/mocks"
	"testing"
)

func TestCourierHandler_Post_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	reqBody := dto.CreateCourierRequest{
		Name:          "John",
		Phone:         "+79779779779",
		Status:        "active",
		TransportType: "on_foot",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	expectedResp := &dto.CreateCourierResponse{
		Id: 10,
	}

	mockService.
		EXPECT().
		CreateCourier(gomock.Any(), &reqBody).
		Return(expectedResp, nil)

	r := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Post(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var decoded dto.CreateCourierResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, 10, decoded.Id)
	require.Equal(t, "courier's profile created successfully", decoded.Message)
}

func TestCourierHandler_Post_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	body := []byte(`invalid json`)

	r := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Post(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInvalidJSON, decoded.Error.Message)
}

func TestCourierHandler_Post_ServiceErrors(t *testing.T) {
	tests := []struct {
		name           string
		req            dto.CreateCourierRequest
		mockResp       *dto.CreateCourierResponse
		mockErr        error
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name: "invalid name error",
			req: dto.CreateCourierRequest{
				Name: "invalid",
			},
			mockResp:       nil,
			mockErr:        service.ErrInvalidName,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidCourierName,
		},
		{
			name: "invalid status error",
			req: dto.CreateCourierRequest{
				Status: "invalid",
			},
			mockResp:       nil,
			mockErr:        service.ErrInvalidStatus,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidCourierStatus,
		},
		{
			name: "invalid phone error",
			req: dto.CreateCourierRequest{
				Phone: "invalid",
			},
			mockResp:       nil,
			mockErr:        service.ErrInvalidPhone,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidCourierPhone,
		},
		{
			name: "courier exists error",
			req: dto.CreateCourierRequest{
				Phone: "exists",
			},
			mockResp:       nil,
			mockErr:        service.ErrCourierExists,
			wantStatusCode: http.StatusConflict,
			wantErrMsg:     server.ErrCourierExists,
		},
		{
			name: "internal error",
			req: dto.CreateCourierRequest{
				Phone: "drop database",
			},
			mockResp:       nil,
			mockErr:        service.ErrInternalError,
			wantStatusCode: http.StatusInternalServerError,
			wantErrMsg:     server.ErrInternalError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_courier.NewMockсourierService(ctrl)
			handler := NewCourierHandler(mockService)

			bodyBytes, _ := json.Marshal(tt.req)

			mockService.
				EXPECT().
				CreateCourier(gomock.Any(), &tt.req).
				Return(tt.mockResp, tt.mockErr)

			r := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.Post(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.wantStatusCode, resp.StatusCode)

			var decoded dto.ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&decoded)
			require.NoError(t, err)

			require.Equal(t, tt.wantErrMsg, decoded.Error.Message)
		})
	}
}

func TestCourierHandler_Get_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Get("/courier/{id}", handler.Get)

	req := dto.GetCourierRequest{
		Id: 10,
	}

	expectedResp := &dto.GetCourierResponse{
		Id:            10,
		Name:          "John",
		Phone:         "+79779779779",
		Status:        "active",
		TransportType: "on_foot",
	}

	mockService.
		EXPECT().
		GetCourier(gomock.Any(), &req).
		Return(expectedResp, nil)

	r := httptest.NewRequest(http.MethodGet, "/courier/10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decoded dto.GetCourierResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, 10, decoded.Id)
	require.Equal(t, "John", decoded.Name)
	require.Equal(t, "+79779779779", decoded.Phone)
	require.Equal(t, "active", decoded.Status)
	require.Equal(t, "on_foot", decoded.TransportType)
}

func TestCourierHandler_Get_InvalidId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Get("/courier/{id}", handler.Get)

	r := httptest.NewRequest(http.MethodGet, "/courier/one", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInvalidCourierId, decoded.Error.Message)
}

func TestCourierHandler_Get_NotFoundError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Get("/courier/{id}", handler.Get)

	req := dto.GetCourierRequest{
		Id: 10,
	}

	mockService.
		EXPECT().
		GetCourier(gomock.Any(), &req).
		Return(nil, service.ErrCourierNotFound)

	r := httptest.NewRequest(http.MethodGet, "/courier/10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrCourierNotFound, decoded.Error.Message)
}

func TestCourierHandler_GetAll_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Get("/couriers", handler.GetAll)

	expectedResp := []dto.GetCourierResponse{
		{
			Id:            10,
			Name:          "John",
			Phone:         "+79779779779",
			Status:        "active",
			TransportType: "on_foot",
		},
		{
			Id:            1,
			Name:          "Rune",
			Phone:         "+79779732779",
			Status:        "busy",
			TransportType: "car",
		},
	}

	mockService.
		EXPECT().
		GetAllCouriers(gomock.Any()).
		Return(expectedResp, nil)

	r := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decoded []dto.GetCourierResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	for i := range decoded {
		require.Equal(t, expectedResp[i].Id, decoded[i].Id)
		require.Equal(t, expectedResp[i].Name, decoded[i].Name)
		require.Equal(t, expectedResp[i].Phone, decoded[i].Phone)
		require.Equal(t, expectedResp[i].Status, decoded[i].Status)
		require.Equal(t, expectedResp[i].TransportType, decoded[i].TransportType)
	}
}

func TestCourierHandler_GetAll_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Get("/couriers", handler.GetAll)

	mockService.
		EXPECT().
		GetAllCouriers(gomock.Any()).
		Return(nil, service.ErrInternalError)

	r := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInternalError, decoded.Error.Message)
}

func TestCourierHandler_Put_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	reqBody := dto.UpdateCourierRequest{
		Name:          "John",
		Phone:         "+79779779779",
		Status:        "active",
		TransportType: "on_foot",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockService.
		EXPECT().
		UpdateCourier(gomock.Any(), &reqBody).
		Return(nil)

	r := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Put(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decoded dto.UpdateCourierResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, "courier's profile updated successfully", decoded.Message)
}

func TestCourierHandler_Put_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	body := []byte(`invalid json`)

	r := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Put(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInvalidJSON, decoded.Error.Message)
}

func TestCourierHandler_Put_ServiceErrors(t *testing.T) {
	tests := []struct {
		name           string
		req            dto.UpdateCourierRequest
		mockResp       *dto.UpdateCourierResponse
		mockErr        error
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name: "invalid name error",
			req: dto.UpdateCourierRequest{
				Name: "invalid",
			},
			mockResp:       nil,
			mockErr:        service.ErrInvalidName,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidCourierName,
		},
		{
			name: "invalid status error",
			req: dto.UpdateCourierRequest{
				Status: "invalid",
			},
			mockResp:       nil,
			mockErr:        service.ErrInvalidStatus,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidCourierStatus,
		},
		{
			name: "invalid phone error",
			req: dto.UpdateCourierRequest{
				Phone: "invalid",
			},
			mockErr:        service.ErrInvalidPhone,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidCourierPhone,
		},
		{
			name: "invalid transport type error",
			req: dto.UpdateCourierRequest{
				TransportType: "invalid",
			},
			mockErr:        service.ErrInvalidTransportType,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     server.ErrInvalidTransportType,
		},
		{
			name: "courier not found",
			req: dto.UpdateCourierRequest{
				Id: -1,
			},
			mockErr:        service.ErrCourierNotFound,
			wantStatusCode: http.StatusNotFound,
			wantErrMsg:     server.ErrCourierNotFound,
		},
		{
			name: "internal error",
			req: dto.UpdateCourierRequest{
				Phone: "drop database",
			},
			mockErr:        service.ErrInternalError,
			wantStatusCode: http.StatusInternalServerError,
			wantErrMsg:     server.ErrInternalError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_courier.NewMockсourierService(ctrl)
			handler := NewCourierHandler(mockService)

			bodyBytes, _ := json.Marshal(tt.req)

			mockService.
				EXPECT().
				UpdateCourier(gomock.Any(), &tt.req).
				Return(tt.mockErr)

			r := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.Put(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.wantStatusCode, resp.StatusCode)

			var decoded dto.ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&decoded)
			require.NoError(t, err)

			require.Equal(t, tt.wantErrMsg, decoded.Error.Message)
		})
	}
}

func TestCourierHandler_Delete_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Delete("/courier/{id}", handler.Delete)

	req := dto.DeleteCourierRequest{
		Id: 10,
	}

	mockService.
		EXPECT().
		DeleteCourier(gomock.Any(), &req).
		Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/courier/10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decoded dto.DeleteCourierResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, "courier's profile deleted successfully", decoded.Message)
}

func TestCourierHandler_Delete_InvalidId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Delete("/courier/{id}", handler.Delete)

	r := httptest.NewRequest(http.MethodDelete, "/courier/one", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInvalidCourierId, decoded.Error.Message)
}

func TestCourierHandler_Delete_NotFoundError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_courier.NewMockсourierService(ctrl)
	handler := NewCourierHandler(mockService)

	router := chi.NewRouter()
	router.Delete("/courier/{id}", handler.Delete)

	req := dto.DeleteCourierRequest{
		Id: 10,
	}

	mockService.
		EXPECT().
		DeleteCourier(gomock.Any(), &req).
		Return(service.ErrCourierNotFound)

	r := httptest.NewRequest(http.MethodDelete, "/courier/10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrCourierNotFound, decoded.Error.Message)
}
