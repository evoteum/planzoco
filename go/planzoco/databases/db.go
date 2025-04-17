package databases

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	DynamoClient *dynamodb.Client
)

func InitDB() error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return err
	}

	DynamoClient = dynamodb.NewFromConfig(cfg)
	return nil
}
