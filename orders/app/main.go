package main

import (
	"context"
	"errors"
	"github.com/halilylm/gommon/rest"
	_orderStream "github.com/halilylm/secondhand/orders/orders/delivery/natstream"
	_productStream "github.com/halilylm/secondhand/orders/product/delivery/natsream"
	_productRepo "github.com/halilylm/secondhand/orders/product/repository/mongodb"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/halilylm/gommon/db"
	"github.com/halilylm/gommon/events/nats"
	"github.com/halilylm/gommon/logger/sugared"
	"github.com/halilylm/gommon/utils"
	_orderHandler "github.com/halilylm/secondhand/orders/orders/delivery/http"
	_orderRepo "github.com/halilylm/secondhand/orders/orders/repository/mongodb"
	_orderUC "github.com/halilylm/secondhand/orders/orders/usecase"
	_productUC "github.com/halilylm/secondhand/orders/product/usecase"
	"github.com/labstack/echo/v4"
)

func main() {
	// parse env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// set the logger
	appLogger := sugared.New(sugared.Options{
		Level:       "info",
		Development: true,
	})
	appLogger.Init()

	// env variables checkpoint
	utils.RequireEnvVariables("MONGO_URI", "APP_PORT", "JWT_KEY", "NATS_URI", "NATS_CLUSTER_ID", "NATS_CLIENT_ID")

	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.NewMongoClient(ctx, os.Getenv("MONGO_URI"))
	if err != nil {
		appLogger.Fatal(err)
	}

	// connect to nats
	streaming, err := nats.New(nats.Options{
		nil,
		appLogger,
		[]string{os.Getenv("NATS_URI")},
		os.Getenv("NATS_CLUSTER_ID"),
		os.Getenv("NATS_CLIENT_ID"),
	})
	if err != nil {
		appLogger.Fatal(err)
	}

	// init collections
	orderCollection := client.Database("orders").Collection("order")
	productCollection := client.Database("orders").Collection("product")

	// init repositories
	orderRepo := _orderRepo.NewOrderRepository(orderCollection)
	productRepo := _productRepo.NewProductRepository(productCollection)

	// init usecases
	orderUC := _orderUC.NewOrder(productRepo, orderRepo, appLogger, streaming)
	ticketUC := _productUC.NewProduct(productRepo, appLogger)

	// set routes
	e := echo.New()
	api := e.Group("/api")
	auth := api.Group("/orders")
	v1 := auth.Group("/v1")

	// 404 handler
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(rest.ErrorResponse(rest.NewNotFoundError()))
	}

	// init handlers
	_orderHandler.NewOrderHandler(v1, orderUC)

	productConsumerGroup := _productStream.NewProductConsumerGroup(streaming, ticketUC, "orders-product-consumer")
	productConsumerGroup.RunConsumers()

	paymentConsumerGroup := _orderStream.NewPaymentConsumerGroup(streaming, orderUC, "orders-payment-consumer")
	paymentConsumerGroup.RunConsumers()

	// start the application
	go func() {
		if err := e.Start(":" + os.Getenv("APP_PORT")); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Fatal(err)
		}
	}()

	// block until program interrupted
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// graceful shutdown
	// don't wait more than 30 seconds
	// to gracefully shut down the server
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		appLogger.Fatal(err)
	}

}
