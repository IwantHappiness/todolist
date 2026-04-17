package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/IwantHappiness/todolist/storage"
	"github.com/IwantHappiness/todolist/task"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var startTime = time.Now()

type MetricsResponse struct {
	MemoryUsageMB uint64 `json:"memory_usage_mb"`
	Uptime        int64  `json:"uptime_seconds"`
	Goroutines    int    `json:"goroutines"`
	CPUCores      int    `json:"cpu_cores"`
}

func newRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/", helloHandler)

	// Tasks
	r.GET("/task", task.GetAllTaskHandler)
	r.GET("/task/:id", task.GetTaskHandler)
	r.POST("/task", task.CreateTaskHandler)
	r.DELETE("/task", task.DeleteAllTaskHandler)
	r.DELETE("/task/:id", task.DeleteTaskHandler)
	r.PATCH("/task/:id", task.CompleteTaskHandler)

	// Metrics
	r.GET("/metrics", metricsHandler)

	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loading .env succesfull")

	log.Println("Server starting")
	router := newRouter()

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	ctx := context.Background()

	dbUrl := os.Getenv("DATABASE_URL")
	_, err := storage.CreateConnection(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	log.Println("Connection db succerfull")

	log.Println("Listen server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stop")
}

func helloHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func metricsHandler(ctx *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	ctx.JSON(http.StatusOK, MetricsResponse{
		Goroutines:    runtime.NumGoroutine(),
		MemoryUsageMB: mem.Alloc / 1024 / 1024,
		CPUCores:      runtime.NumCPU(),
		Uptime:        int64(time.Since(startTime).Seconds()),
	})
}
