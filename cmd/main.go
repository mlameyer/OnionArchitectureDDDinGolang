package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/infrastructure/awsservice"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Name     string `json:"DB_NAME"`
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
}

func main() {
	// Retrieve database configuration from AWS Secrets Manager
	secretsManager := &awsservice.AwsSecretsManager{}
	secretValue, err := secretsManager.GetSecretValue(context.Background(), os.Getenv("AWS_SECRET_NAME"))
	if err != nil {
		log.Fatalf("failed to retrieve secret value: %v", err)
	}

	var dbConfig DBConfig
	if err := json.Unmarshal([]byte(secretValue), &dbConfig); err != nil {
		log.Fatalf("failed to unmarshal secret string: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	log.Printf("Connected to the database successfully: %v", db)

	// Apply migrations
	err = db.AutoMigrate(&Order{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully.")

	// Your application logic here
}
