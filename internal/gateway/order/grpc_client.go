package order

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"service-order-avito/api/order"
	"time"
)

type orderServiceRPCClient interface {
	GetOrders(ctx context.Context, in *order.GetOrdersRequest, opts ...grpc.CallOption) (*order.GetOrdersResponse, error)
	GetOrderById(ctx context.Context, in *order.GetOrderByIdRequest, opts ...grpc.CallOption) (*order.GetOrderByIdResponse, error)
}

type orderGateway struct {
	client orderServiceRPCClient
}

func NewOrderGateway(c orderServiceRPCClient) *orderGateway {
	return &orderGateway{client: c}
}

func (og *orderGateway) GetOrderIdsFrom(ctx context.Context, from time.Time) ([]string, error) {
	resp, err := og.client.GetOrders(
		ctx,
		&order.GetOrdersRequest{From: timestamppb.New(from)},
	)
	if err != nil {
		return nil, err
	}

	orderIds := make([]string, len(resp.Orders))
	for i := range resp.Orders {
		orderIds[i] = resp.Orders[i].OrderNumber
	}

	return orderIds, err
}

func (og *orderGateway) GetOrderStatusById(ctx context.Context, id string) (string, error) {
	resp, err := og.client.GetOrderById(
		ctx,
		&order.GetOrderByIdRequest{Id: id},
	)
	if err != nil {
		return "", err
	}

	return resp.Order.GetStatus(), nil
}
