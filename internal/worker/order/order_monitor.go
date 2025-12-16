package order

import (
	"context"
	"log/slog"
	"service-order-avito/internal/domain/dto"
	"time"
)

type deliveryService interface {
	AssignDelivery(context.Context, *dto.AssignDeliveryRequest) (*dto.AssignDeliveryResponse, error)
}

type orderServiceRPCClient interface {
	GetOrderIdsFrom(context.Context, time.Time) ([]string, error)
}

type orderMonitorWorker struct {
	interval   time.Duration
	log        *slog.Logger
	delService deliveryService
	RPCClient  orderServiceRPCClient
}

func NewOrderMonitorWorker(interval time.Duration, log *slog.Logger, delService deliveryService, client orderServiceRPCClient) *orderMonitorWorker {
	return &orderMonitorWorker{
		interval:   interval,
		log:        log,
		delService: delService,
		RPCClient:  client,
	}
}

func (w *orderMonitorWorker) Start(ctx context.Context) {
	timeCursor := time.Now()
	ticker := time.NewTicker(w.interval)

	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			w.log.Info("delivery monitor worker gracefully stopped")
			return
		case <-ticker.C:
			orderIds, err := w.RPCClient.GetOrderIdsFrom(ctx, timeCursor)
			if err != nil {
				w.log.Error(err.Error())
				continue
			}

			for _, orderId := range orderIds {
				req, err := w.delService.AssignDelivery(
					ctx,
					&dto.AssignDeliveryRequest{OrderId: orderId},
				)
				if err != nil {
					w.log.Error(err.Error())
					continue
				}

				w.log.Info("order: " + req.OrderId + " was successfully assigned")
			}

			timeCursor.Add(w.interval)
		}

	}
}
