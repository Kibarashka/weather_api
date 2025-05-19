package main

import (
	"fmt"
	"log"
	"weather/project/client"
	"weather/project/config"
	"weather/project/handler"
	"weather/project/repository"
	"weather/project/server"
	"weather/project/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("FATAL: Could not load .env config: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("FATAL: Could not initialize database: %v", err)
	}
	log.Println("Database initialized successfully.")

	if err := repository.MigrateDB(db); err != nil {
		log.Fatalf("FATAL: Could not migrate database: %v", err)
	}
	log.Println("Database migration completed successfully.")

	weatherAPIClient := client.NewWeatherAPIClient(cfg)

	subscriptionRepo := repository.NewSubscriptionRepository(db)

	tokenSvc := service.NewTokenService()
	emailSvc := service.NewEmailService(cfg) // Pass cfg for AppBaseURL etc.
	subscriptionSvc := service.NewSubscriptionService(subscriptionRepo, tokenSvc, emailSvc)
	weatherSvc := service.NewWeatherService(weatherAPIClient)

	weatherHdlr := handler.NewWeatherHandler(weatherSvc)
	subscriptionHdlr := handler.NewSubscriptionHandler(subscriptionSvc)
	log.Println("Dependencies initialized.")

	router := server.SetupRouter(weatherHdlr, subscriptionHdlr)
	log.Println("HTTP router setup complete.")

	appAddress := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting Weather API server on %s", appAddress)
	log.Printf("API Documentation available at http://localhost:%s/swagger.yaml", cfg.AppPort)

	if err := router.Run(appAddress); err != nil {
		log.Fatalf("FATAL: Could not start server on %s: %v", appAddress, err)
	}
}
