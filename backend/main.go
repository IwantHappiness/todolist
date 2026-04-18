package main

import (
	"context"
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
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Loading .env succesful")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v:", err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected db succerful")

	repo := repository.NewTaskPgRepository(conn)
	scv := service.NewTaskService(repo)
	handler := handler.NewTaskHandler(scv)

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	handler.RegisterRouter(router)

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Println("Server starting")

	go func() {
		log.Println("Server running on port " + port)
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
