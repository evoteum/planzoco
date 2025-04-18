package databases

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DB *dynamodb.Client

func InitDB(ctx context.Context) error {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		return fmt.Errorf("DYNAMODB_TABLE not set")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	DB = dynamodb.NewFromConfig(cfg)

	_, err = DB.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &tableName,
	})
	if err != nil {
		return fmt.Errorf("could not describe table %s: %w", tableName, err)
	}

	return nil
}