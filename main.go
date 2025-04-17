package main

import (
	"accounts-service/config"
	"accounts-service/handlers"
	"accounts-service/repositories"
	"accounts-service/usecases"
	"accounts-service/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Parse command line arguments
	args := utils.ParseArguments()

	// Load configuration
	cfg, err := config.LoadConfig(args.ConfigPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize for logger
	logger := utils.NewLogger(cfg.LogLevel)

	// Initialize database connection
	db, err := config.NewDatabaseConnection(cfg)
	if err != nil {
		logger.Critical("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	accountRepo := repositories.NewAccountRepository(db, logger)
	mutationRepo := repositories.NewMutationRepository(db, logger)

	// Initialize usecase
	accountUsecase := usecases.NewAccountUsecase(accountRepo, mutationRepo, logger)

	// Initialize handler
	accountHandler := handlers.NewAccountHandler(accountUsecase, logger)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// Routes
	api := e.Group("/api/account")

	api.POST("/daftar", accountHandler.CreateAccount)
	api.POST("/tabung", accountHandler.Credit)
	api.POST("/tarik", accountHandler.Debit)
	api.GET("/saldo/:no_rekening", accountHandler.GetSaldo)

	// Start server
	go func() {
		if err := e.Start(":" + cfg.AppPort); err != nil {
			logger.Error("Shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error: %v", err)
	}
}
