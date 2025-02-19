package main

import (
	configs "cakestore/config"
	controller "cakestore/internal/delivery/http"
	"cakestore/internal/delivery/http/route"
	"cakestore/internal/repository"
	"cakestore/internal/usecase"
	"cakestore/utils"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	logger := utils.NewLogger()

	cfg := configs.LoadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	app := fiber.New()

	// repo
	cakeRepository := repository.NewCakeRepository(db, logger)
	// usecase
	cakeUseCase := usecase.NewCakeUseCase(cakeRepository, logger)
	// controller
	cakeController := controller.NewCakeController(cakeUseCase, logger)

	routeConfig := route.RouteConfig{
		App:            app,
		CakeController: cakeController,
	}
	routeConfig.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
