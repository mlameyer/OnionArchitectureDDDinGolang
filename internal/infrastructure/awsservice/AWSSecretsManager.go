package awsservice

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecretsManager struct{}

func (a *AwsSecretsManager) GetSecretValue(ctx context.Context, secretId string) (string, error) {
	// Load AWS config with the default credentials provider
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	// Retrieve secret value
	result, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret value, %v", err)
	}

	// Return the secret string
	return *result.SecretString, nil
}
