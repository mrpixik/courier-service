package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"service-order-avito/api/order"
	"service-order-avito/internal/config"
	order2 "service-order-avito/internal/gateway/order"
	"service-order-avito/internal/http/server"
	courier2 "service-order-avito/internal/http/server/handlers/courier"
	delivery2 "service-order-avito/internal/http/server/handlers/delivery"
	"service-order-avito/internal/repository/postgres"
	"service-order-avito/internal/service/courier"
	"service-order-avito/internal/service/delivery"
	delivery_worker "service-order-avito/internal/worker/delivery"
	order_worker "service-order-avito/internal/worker/order"
	"service-order-avito/pkg/logger"
	"syscall"
	"time"
)

func main() {
	// CONFIG
	cfg := config.MustLoad()

	// LOGGER
	log := logger.MustInit(cfg.Env)
	log.Info("logger initialized")

	// GS context
	ctxApp, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// context для бд, отменяется после завершения сервера
	ctxDB, cancelDB := context.WithCancel(context.Background())
	defer cancelDB()

	// DB connection
	pool, err := postgres.ConnectPostgres(ctxDB, cfg.Postgres, cfg.Env)
	if err != nil {
		log.Error("connect database " + err.Error())
		os.Exit(1)
	}
	defer func() {
		pool.Close()
		log.Info("connection with database closed")
	}()

	// Repository Lay
	transactionManager := postgres.NewTransactionManagerPostgres(pool)
	courierRepository := postgres.NewCourierRepositoryPostgres(pool)
	deliveryRepository := postgres.NewDeliveryRepositoryPostgres(pool)
	log.Info("repository lay is initialized")

	// Service lay
	courierService := courier.NewCourierService(transactionManager, courierRepository)
	deliveryService := delivery.NewDeliveryService(transactionManager, courierRepository, deliveryRepository)
	log.Info("service lay is initialized")

	// Workers
	deliveryMonitorWorker := delivery_worker.NewDeliveryMonitorWorker(cfg.DeliveryWorkerTickInterval, log, deliveryService)
	go deliveryMonitorWorker.Start(ctxApp)
	log.Info("delivery monitor worker is started")

	connRPC, err := grpc.NewClient(cfg.GRPC.OrderServiceDSN, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("unable to connect to order-service")
		os.Exit(1)
	}
	orderServiceClient := order.NewOrdersServiceClient(connRPC)
	orderGateway := order2.NewOrderGateway(orderServiceClient)

	orderServiceMonitorWorker := order_worker.NewOrderMonitorWorker(
		time.Second*5,
		log,
		deliveryService,
		orderGateway,
	)
	go orderServiceMonitorWorker.Start(ctxApp)
	log.Info("order-service monitor worker is started")

	// Controller lay
	courierHandler := courier2.NewCourierHandler(courierService)
	deliveryHandler := delivery2.NewDeliveryHandler(deliveryService)
	log.Info("controller lay is initialized")

	// ROUTER & SERVER
	r := server.InitRouter(cfg.HTTP, log, courierHandler, deliveryHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: r,
		BaseContext: func(net.Listener) context.Context {
			return ctxApp
		},
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.ShutdownTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server start up: %s", err.Error())
		}
	}()
	log.Info("listening on: " + cfg.HTTP.Port)

	gracefulShutdown(ctxApp, cfg.HTTP, log, srv, cancelDB)
}

func gracefulShutdown(ctxApp context.Context, cfg config.HTTPServer, log *slog.Logger, srv *http.Server, cancelDB context.CancelFunc) {
	<-ctxApp.Done()
	log.Info("shutdown signal received. starting graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown: %s", err.Error())
	} else {
		log.Info("server gracefully stopped")
	}

	cancelDB()
}
