package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/handler/http/server/handler/delivery/mocks"
	"testing"
	"time"
)

func TestCourierHandler_PostAssign_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_delivery.NewMockdeliveryService(ctrl)
	handler := NewDeliveryHandler(mockService)

	reqBody := dto.AssignDeliveryRequest{
		OrderId: "SWAG--SWAG--SWAG",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	expectedResp := &dto.AssignDeliveryResponse{
		CourierId:        1,
		OrderId:          "SWAG--SWAG--SWAG",
		TransportType:    "helicopter",
		DeliveryDeadline: time.Date(1997, time.August, 29, 0, 0, 0, 0, time.UTC),
	}

	mockService.
		EXPECT().
		AssignDelivery(gomock.Any(), &reqBody).
		Return(expectedResp, nil)

	r := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.PostAssign(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var decoded dto.AssignDeliveryResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, expectedResp.CourierId, decoded.CourierId)
	require.Equal(t, expectedResp.OrderId, decoded.OrderId)
	require.Equal(t, expectedResp.TransportType, decoded.TransportType)
	require.Equal(t, expectedResp.DeliveryDeadline, decoded.DeliveryDeadline)
}

func TestCourierHandler_PostAssign_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_delivery.NewMockdeliveryService(ctrl)
	handler := NewDeliveryHandler(mockService)

	body := []byte(`invalid json`)

	r := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostAssign(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInvalidJSON, decoded.Error.Message)
}

func TestCourierHandler_PostAssign_ServiceErrors(t *testing.T) {
	tests := []struct {
		name           string
		req            dto.AssignDeliveryRequest
		mockResp       *dto.AssignDeliveryResponse
		mockErr        error
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name: "no available couriers",
			req: dto.AssignDeliveryRequest{
				OrderId: "777-SWAG-777",
			},
			mockResp:       nil,
			mockErr:        service.ErrNoAvailableCouriers,
			wantStatusCode: http.StatusConflict,
			wantErrMsg:     server.ErrNoAvailableCouriers,
		},
		{
			name: "delivery already exists",
			req: dto.AssignDeliveryRequest{
				OrderId: "777-SWAG-777",
			},
			mockResp:       nil,
			mockErr:        service.ErrDeliveryExists,
			wantStatusCode: http.StatusConflict,
			wantErrMsg:     server.ErrDeliveryExists,
		},
		{
			name: "internal error",
			req: dto.AssignDeliveryRequest{
				OrderId: "777-SWAG-777",
			},
			mockResp:       nil,
			mockErr:        service.ErrInternalError,
			wantStatusCode: http.StatusInternalServerError,
			wantErrMsg:     server.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_delivery.NewMockdeliveryService(ctrl)
			handler := NewDeliveryHandler(mockService)

			bodyBytes, _ := json.Marshal(tt.req)

			mockService.
				EXPECT().
				AssignDelivery(gomock.Any(), &tt.req).
				Return(tt.mockResp, tt.mockErr)

			r := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.PostAssign(w, r)

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

func TestCourierHandler_PostUnassign_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_delivery.NewMockdeliveryService(ctrl)
	handler := NewDeliveryHandler(mockService)

	reqBody := dto.UnassignDeliveryRequest{
		OrderId: "SWAG--SWAG--SWAG",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	expectedResp := &dto.UnassignDeliveryResponse{
		OrderId:   "SWAG--SWAG--SWAG",
		Status:    "unassigned",
		CourierId: 1,
	}

	mockService.
		EXPECT().
		UnassignDelivery(gomock.Any(), &reqBody).
		Return(expectedResp, nil)

	r := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.PostUnassign(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decoded dto.UnassignDeliveryResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, expectedResp.CourierId, decoded.CourierId)
	require.Equal(t, expectedResp.OrderId, decoded.OrderId)
	require.Equal(t, expectedResp.Status, decoded.Status)
}

func TestCourierHandler_PostUnassign_InvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_delivery.NewMockdeliveryService(ctrl)
	handler := NewDeliveryHandler(mockService)

	body := []byte(`invalid json`)

	r := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.PostUnassign(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var decoded dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&decoded)
	require.NoError(t, err)

	require.Equal(t, server.ErrInvalidJSON, decoded.Error.Message)
}

func TestCourierHandler_PostUnassign_ServiceErrors(t *testing.T) {
	tests := []struct {
		name           string
		req            dto.UnassignDeliveryRequest
		mockResp       *dto.UnassignDeliveryResponse
		mockErr        error
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name: "delivery not found",
			req: dto.UnassignDeliveryRequest{
				OrderId: "777-SWAG-777",
			},
			mockResp:       nil,
			mockErr:        service.ErrDeliveryNotFound,
			wantStatusCode: http.StatusNotFound,
			wantErrMsg:     server.ErrDeliveryNotFound,
		},
		{
			name: "internal error",
			req: dto.UnassignDeliveryRequest{
				OrderId: "777-SWAG-777",
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

			mockService := mock_delivery.NewMockdeliveryService(ctrl)
			handler := NewDeliveryHandler(mockService)

			bodyBytes, _ := json.Marshal(tt.req)

			mockService.
				EXPECT().
				UnassignDelivery(gomock.Any(), &tt.req).
				Return(tt.mockResp, tt.mockErr)

			r := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.PostUnassign(w, r)

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
