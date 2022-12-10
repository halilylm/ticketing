package main

import (
	"context"
	"errors"
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
	ticketCollection := client.Database("ticket").Collection("tickets")

	// init repositories
	ticketRepo := mongodb.NewTicketRepository(ticketCollection)

	// init nats streaming
	streaming, err := nats.New(nats.Options{
		nil,
		appLogger,
		[]string{"nats://localhost:4222"},
		"test-cluster",
		"client_id",
	})
	if err != nil {
		appLogger.Fatal(err)
	}

	// init usecases
	ticketUC := usecase.NewTicket(ticketRepo, appLogger, streaming)

	// set routes
	e := echo.New()
	api := e.Group("/api")
	ticket := api.Group("/tickets")
	v1 := ticket.Group("/v1")

	// init handlers
	_ticketHandler.NewTicketHandler(v1, ticketUC)

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
