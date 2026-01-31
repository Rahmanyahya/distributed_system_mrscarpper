package main

import (
	"distributed_system/internal/config"
	"distributed_system/internal/infrastructure/database"
	"distributed_system/seeds"
	"fmt"
	"log"
	"os"
)

func main() {
	// Load config
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	fmt.Println("==========================================")
	fmt.Println("     ADMIN SEEDING UTILITIES")
	fmt.Println("==========================================")
	fmt.Println()

	// Run default admin seed
	fmt.Println("Running admin seed with default credentials...")
	fmt.Println()

	if err := seeds.RunAdminSeed(db); err != nil {
		log.Fatalf("Failed to seed admin: %v\n", err)
	}

	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("     SEEDING COMPLETED SUCCESSFULLY!")
	fmt.Println("==========================================")
	fmt.Println()
	fmt.Println("📝 Default Admin Credentials:")
	fmt.Println("   Email: admin@distributed-system.com")
	fmt.Println("   Password: Admin123!@#")
	fmt.Println()
	fmt.Println("⚠️  IMPORTANT: Please change the password after first login!")
	fmt.Println()
}
