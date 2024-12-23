package main

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/domain/models"

	"order-service/internal/application/handlers"
	"order-service/internal/domain/services"
	"order-service/internal/infrastructure/awsservice"
	"order-service/internal/infrastructure/logging"
	"order-service/internal/infrastructure/persistence"
	"order-service/internal/infrastructure/tracing"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title Order Service API
// @version 1.0
// @description This is an API for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080

type DBConfig struct {
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Name     string `json:"DB_NAME"`
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
}

func main() {
	// Initialize logging
	logging.InitLogger()

	// Initialize tracing
	if err := tracing.InitTracer(); err != nil {
		logging.Logger.Error().Msgf("failed to initialize tracer: %v", err)
		return
	}
	defer func() {
		if err := tracing.ShutdownTracer(context.Background()); err != nil {
			logging.Logger.Error().Msgf("failed to shutdown tracer: %v", err)
		}
	}()

	// Retrieve database configuration from AWS Secrets Manager
	secretsManager := &awsservice.AwsSecretsManager{}
	secretValue, err := secretsManager.GetSecretValue(context.Background(), os.Getenv("AWS_SECRET_NAME"))
	if err != nil {
		logging.Logger.Error().Msgf("failed to retrieve secret value: %v", err)
	}

	var dbConfig DBConfig
	if err := json.Unmarshal([]byte(secretValue), &dbConfig); err != nil {
		logging.Logger.Error().Msgf("failed to unmarshal secret string: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logging.Logger.Error().Msgf("could not connect to the database: %v", err)
	}

	logging.Logger.Info().Msgf("Connected to the database successfully: %v", db)

	// Apply migrations
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		logging.Logger.Error().Msgf("failed to migrate database: %v", err)
	}

	err = db.AutoMigrate(&models.OrderItem{})
	if err != nil {
		logging.Logger.Error().Msgf("failed to migrate database: %v", err)
	}

	logging.Logger.Info().Msg("Database migration completed successfully.")

	// Set up repositories
	orderRepo := persistence.NewGormOrderRepository(db)

	// Set up event publisher
	eventPublisher := &services.LoggerEventPublisher{}

	// Set up services
	orderService := services.NewOrderService(orderRepo, eventPublisher)

	// Set up Fiber and API handlers
	app := fiber.New()
	handlers.NewOrderHandler(app, orderService)

	// Serve Swagger UI
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}))

	// Start the server
	logging.Logger.Info().Msg("Starting server on port 8080")
	if err := app.Listen(":8080"); err != nil {
		logging.Logger.Error().Msgf("failed to start server: %v", err)
	}
}
