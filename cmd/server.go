package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chand-magar/SolidBaseGoStructure/internal/config"
	database "github.com/chand-magar/SolidBaseGoStructure/internal/database"
	"github.com/chand-magar/SolidBaseGoStructure/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.MustLoad()

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Environment Variables
	dbHost := os.Getenv("_DATABASE_HOST_")
	dbUser := os.Getenv("_DATABASE_USER_")
	dbPass := "M@gAr!~t0rE!#2025" //os.Getenv("_DATABASE_PASSWORD_")
	dbName := os.Getenv("_DATABASE_NAME_")
	dbPort := os.Getenv("_DATABASE_PORT_")

	if dbPort == "" {
		dbPort = "5433" // Default PostgreSQL port
	}

	if dbHost == "" || dbUser == "" || dbPass == "" || dbName == "" {
		fmt.Println("Environment variables not set properly.")
		os.Exit(1)
	}

	// Construct DSN from environment variables
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		dbHost, dbUser, dbPass, dbName, dbPort)

	// Initialize DB
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	switch cfg.GinMode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		gin.SetMode(cfg.GinMode)
	default:
		log.Fatalf("Invalid GIN_MODE: %s", cfg.GinMode)
	}

	// Initialize Gin router with DB
	ginRouter := router.AllRouter(db)

	// Create http.Server with ginRouter as Handler
	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: ginRouter,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
