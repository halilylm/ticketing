package main

import (
	"context"
	"errors"
	"github.com/halilylm/gommon/db"
	"github.com/halilylm/gommon/logger/sugared"
	"github.com/halilylm/gommon/utils"
	_orderHandler "github.com/halilylm/ticketing/orders/orders/delivery/http"
	_orderRepo "github.com/halilylm/ticketing/orders/orders/repository/mongodb"
	"github.com/halilylm/ticketing/orders/orders/usecase"
	_ticketRepo "github.com/halilylm/ticketing/orders/ticket/repository/mongodb"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// set the logger
	appLogger := sugared.New(sugared.Options{
		Level:       "info",
		Development: true,
	})
	appLogger.Init()

	// env variables checkpoint
	utils.RequireEnvVariables("MONGO_URI", "APP_PORT", "JWT_KEY")

	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.NewMongoClient(ctx, os.Getenv("MONGO_URI"))
	if err != nil {
		appLogger.Fatal(err)
	}

	// init collections
	orderCollection := client.Database("orders").Collection("orders")
	ticketCollection := client.Database("orders").Collection("tickets")
	appLogger.Info(ticketCollection.CountDocuments(ctx, bson.M{}))

	// init repositories
	orderRepo := _orderRepo.NewOrderRepository(orderCollection)
	ticketRepo := _ticketRepo.NewTicketRepository(ticketCollection)

	// init usecases
	orderUC := usecase.NewOrder(ticketRepo, orderRepo, appLogger)

	// set routes
	e := echo.New()
	api := e.Group("/api")
	auth := api.Group("/orders")
	v1 := auth.Group("/v1")

	// init handlers
	_orderHandler.NewOrderHandler(v1, orderUC)

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
