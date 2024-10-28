package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Name     string `json:"DB_NAME"`
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
}

type Order struct {
	OrderId string `json:`
}

func main() {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	// Retrieve secret value
	secretName := os.Getenv("AWS_SECRET_NAME")
	fmt.Printf("my secret name is: %s\n", secretName)
	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatalf("failed to retrieve secret value, %v", err)
	}

	// Parse the secret JSON
	secretString := *result.SecretString
	var dbConfig DBConfig
	if err := json.Unmarshal([]byte(secretString), &dbConfig); err != nil {
		log.Fatalf("failed to unmarshal secret string, %v", err)
	}

	dsn := "host=" + dbConfig.Host + " user=" + dbConfig.User + " password=" + dbConfig.Password + " dbname=" + dbConfig.Name + " port=" + dbConfig.Port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	log.Println("Connected to the database successfully: %v", db)

	// Apply migrations using golang-migrate
	cmd := exec.Command("migrate", "-path", "sql/migrations", "-database", dsn, "up")
	if err := cmd.Run(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Your application logic here
}
