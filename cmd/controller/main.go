package main

import (
	"distributed_system/internal/config"
	"distributed_system/internal/delivery/http/handler"
	"distributed_system/internal/delivery/http/middleware"
	"distributed_system/internal/infrastructure/cache"
	"distributed_system/internal/infrastructure/database"
	"distributed_system/internal/infrastructure/redis"
	"distributed_system/internal/repository/admin"
	"distributed_system/internal/repository/agents"
	configRepo "distributed_system/internal/repository/config"
	adminUC "distributed_system/internal/usecase/admin"
	agentUC "distributed_system/internal/usecase/agents"
	configUC "distributed_system/internal/usecase/config"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	servicePort := cfg.Server.Controller.Port

	db := initDatabase(cfg)
	defer db.Close()

	redisClient := initRedis(cfg)

	configCache := cache.NewConfigCache(redisClient)
	configRepository := configRepo.NewCOnfigRepository(db.DB, configCache)
	agentsRepository := agents.NewAgentRepository(db.DB)
	adminRepository := admin.NewAdminRepository(db.DB)

	configUsecase := configUC.NewConfigUsecase(configRepository, agentsRepository, cfg, configCache)
	agentsUsecase := agentUC.NewAgentUsecase(agentsRepository, cfg)
	adminUsecase := adminUC.NewAdminUsecase(adminRepository, cfg)

	configHandler := handler.NewConfigHandler(configUsecase)
	agentHandler := handler.NewAgentsHandler(agentsUsecase)
	adminHandler := handler.NewAdminHandler(adminUsecase)

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", adminHandler.Login)

	groupConfig := r.Group("/config")
	{
		admin := groupConfig.Group("/admin")
		{
			admin.Use(middleware.AdminValidation(cfg))
			admin.GET("", configHandler.GetLatestConfigAdmin)
			admin.PUT("", configHandler.Update)
			admin.POST("", configHandler.Create)
		}

		agent := groupConfig.Group("/agent") 
		{
			agent.Use(middleware.InternalGetConfigVaidation(cfg))
			agent.GET("", configHandler.GetLatestConfigModel)
		}

	}

	groupAgent := r.Group("/agent")
	{
		register := groupAgent.Group("/register")
		{
			register.Use(middleware.ValidationRegistrationAgent(cfg))
			register.POST("", agentHandler.Register)
		}

		admin := groupAgent.Group("/admin")
		{
			admin.Use(middleware.AdminValidation(cfg))
			admin.GET("", agentHandler.GenerateRegistrationConfifg)
		}
	}

	r.Run(fmt.Sprintf(":%d", servicePort))
}

func initDatabase(cfg *config.Config) *database.Database {
	db, err := database.New(&cfg.Database)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	return db
}

func initRedis(cfg *config.Config) *redis.Client {
	redisClient, err := redis.New(&cfg.Redis)
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	return redisClient
}