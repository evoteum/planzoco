package databases

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	TableName = "planzoco"
	DefaultRegion = "eu-west-2"
)

var (
	DynamoClient *dynamodb.Client
)

func GetTableName() string {
	if name := os.Getenv("DYNAMODB_TABLE"); name != "" {
		return name
	}
	return TableName
}

func GetRegion() string {
	if region := os.Getenv("AWS_REGION"); region != "" {
		return region
	}
	return DefaultRegion
}

func InitDB() error {
	region := GetRegion()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Printf("Error: unable to load SDK config, %v", err)
		return err
	}

	DynamoClient = dynamodb.NewFromConfig(cfg)
	log.Printf("DynamoDB client initialized, using table: %s in region: %s", GetTableName(), region)

	return nil
}
