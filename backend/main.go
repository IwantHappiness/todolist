package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/IwantHappiness/TodoList/task"
	"github.com/gin-gonic/gin"
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
	log.Println("Server starting")

	router := newRouter()

	port := getEnv("PORT", "8080")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server stop")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
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
