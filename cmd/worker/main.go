package main

import (
	"context"
	"crypto/tls"
	"distributed_system/internal/config"
	"distributed_system/internal/delivery/http/handler"
	"distributed_system/internal/delivery/http/middleware"
	"distributed_system/internal/usecase/worker"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get config path from env or use default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	workerCfg, err := config.LoadWorkerConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load worker config: %v", err)
	}

	log.Println("============================================================")
	log.Println("[Worker] Starting...")
	log.Printf("[Worker] Port: %d", workerCfg.Server.Port)
	log.Printf("[Worker] Waiting for Agent to push initial config...")
	log.Println("============================================================")

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		// Allow insecure HTTPS for development (self-signed certificates)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	workerUsecase := worker.NewWorkerUsecase(httpClient)
	workerHandler := handler.NewWorkerHandler(workerUsecase)

	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/hit", workerHandler.Hit)
	
	privateGroup := r.Group("/private")
	{
		privateGroup.Use(middleware.ValidationAgentWorker(workerCfg))
		privateGroup.POST("/config", workerHandler.UpdateConfig)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "worker",
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", workerCfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Printf("[Worker] Server started on port %d...", workerCfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Worker server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("[Worker] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Worker forced to shutdown: %v", err)
	}

	log.Println("[Worker] Server stopped.")
}
