package main

import (
	"context"
	"errors"
	"github.com/halilylm/gommon/db"
	"github.com/halilylm/gommon/logger/sugared"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/gommon/utils"
	_authHandler "github.com/halilylm/ticketing/auth/auth/delivery/http"
	"github.com/halilylm/ticketing/auth/auth/repository/mongodb"
	"github.com/halilylm/ticketing/auth/auth/usecase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
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
	utils.RequireEnvVariables("MONGO_URI", "APP_PORT", "JWT_KEY")

	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.NewMongoClient(ctx, os.Getenv("MONGO_URI"))
	if err != nil {
		appLogger.Fatal(err)
	}

	// init collections
	userCollection := client.Database("auth").Collection("users")

	// init repositories
	userRepo := mongodb.NewUserRepository(userCollection)

	// init usecases
	authUC := usecase.NewAuth(userRepo, appLogger)

	// set routes
	e := echo.New()
	api := e.Group("/api")
	auth := api.Group("/auth")
	v1 := auth.Group("/v1")

	// 404 handler
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(rest.ErrorResponse(rest.NewNotFoundError()))
	}

	// init handlers
	_authHandler.NewAuthHandler(v1, authUC)

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
