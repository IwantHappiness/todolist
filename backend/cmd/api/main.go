package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IwantHappiness/todolist/internal/handler"
	"github.com/IwantHappiness/todolist/internal/repository"
	"github.com/IwantHappiness/todolist/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	cfg := loadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := pgx.Connect(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v:", err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected db succerful")

	repo := repository.NewTaskPgRepository(conn)
	scv := service.NewTaskService(repo)
	handler := handler.NewTaskHandler(scv)

	gin.SetMode(cfg.GIN_MODE)
	router := gin.Default()
	handler.RegisterRouter(router)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	log.Println("Server starting")

	go func() {
		log.Println("Server running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}

type config struct {
	HTTPAddr    string
	DatabaseDSN string
	GIN_MODE    string
}

func loadConfig() config {
	cfg := config{
		HTTPAddr:    envOrDefault("HTTP_ADDR", ":8080"),
		DatabaseDSN: envOrDefault("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/taskservice?sslmode=disable"),
		GIN_MODE:    envOrDefault("GIN_MODE", "release"),
	}

	if cfg.DatabaseDSN == "" {
		panic(fmt.Errorf("DATABASE_DSN is required"))
	}

	return cfg
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
