package awsservice

import "context"

type SecretsManager interface {
	GetSecretValue(ctx context.Context, secretId string) (string, error)
}
