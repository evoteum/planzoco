package databases

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	// TableName is the single table name for all entities
	TableName = "planzoco"
	// DefaultRegion is used if no region is specified
	DefaultRegion = "eu-west-2"
)

var (
	DynamoClient *dynamodb.Client
)

// GetTableName returns the table name based on environment variables or defaults
func GetTableName() string {
	if name := os.Getenv("DYNAMODB_TABLE"); name != "" {
		return name
	}
	return TableName
}

// GetRegion returns the AWS region to use
func GetRegion() string {
	if region := os.Getenv("AWS_REGION"); region != "" {
		return region
	}
	return DefaultRegion
}

// InitDB initializes the DynamoDB client
func InitDB() error {
	// Get region from environment or use default
	region := GetRegion()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return err
	}

	// Initialize DynamoDB client
	DynamoClient = dynamodb.NewFromConfig(cfg)
	log.Printf("DynamoDB client initialized, using table: %s in region: %s", GetTableName(), region)

	return nil
}
