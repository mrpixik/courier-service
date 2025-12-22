package delivery

import (
	"context"
	"fmt"
	"service-order-avito/internal/adapters/logger"
	"time"
)

type deliveryService interface {
	UnassignAllCompleted(context.Context) (int, error)
}

type deliveryMonitorWorker struct {
	interval   time.Duration
	log        logger.LoggerAdapter
	delService deliveryService
}

func NewDeliveryMonitorWorker(interval time.Duration, log logger.LoggerAdapter, delService deliveryService) *deliveryMonitorWorker {
	return &deliveryMonitorWorker{interval: interval, log: log, delService: delService}
}

func (w *deliveryMonitorWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			w.log.Info("delivery monitor worker gracefully stopped")
			return
		case <-ticker.C:
			totalUnassigned, err := w.delService.UnassignAllCompleted(ctx)
			if err != nil {
				w.log.Error(err.Error())
				continue
			}

			if totalUnassigned > 0 {
				w.log.Info(fmt.Sprintf("unassigned %d deliveries", totalUnassigned))
			}
		}

	}
}
