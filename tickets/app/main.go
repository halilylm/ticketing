package main

import (
	"context"
	"errors"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/ticketing/tickets/ticket/delivery/natstream"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/halilylm/gommon/db"
	"github.com/halilylm/gommon/events/nats"
	"github.com/halilylm/gommon/logger/sugared"
	"github.com/halilylm/gommon/utils"
	_ticketHandler "github.com/halilylm/ticketing/tickets/ticket/delivery/http"
	"github.com/halilylm/ticketing/tickets/ticket/repository/mongodb"
	"github.com/halilylm/ticketing/tickets/ticket/usecase"
	"github.com/labstack/echo/v4"

	"os"
	"os/signal"
	"time"
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
	ticketCollection := client.Database("ticket").Collection("tickets")

	// init repositories
	ticketRepo := mongodb.NewTicketRepository(ticketCollection)

	// init usecases
	ticketUC := usecase.NewTicket(ticketRepo, appLogger, streaming)

	// set routes
	e := echo.New()
	api := e.Group("/api")
	ticket := api.Group("/tickets")
	v1 := ticket.Group("/v1")

	// 404 handler
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(rest.ErrorResponse(rest.NewNotFoundError()))
	}

	// init handlers
	_ticketHandler.NewTicketHandler(v1, ticketUC)

	orderConsumerGroup := natstream.NewOrderConsumerGroup(streaming, ticketUC, "ticket_order_consumer")
	orderConsumerGroup.RunConsumers()

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
