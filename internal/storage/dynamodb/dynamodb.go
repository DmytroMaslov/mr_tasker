package dynamodb

import (
	"context"
	"fmt"
	"mr-tasker/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -source=dynamodb.go -destination=mocks/dynamodb_mock.go
type DynamoDbClient interface {
	AddRow(context.Context, any) (*dynamodb.PutItemOutput, error)
	ReadRow(context.Context, map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error)
	UpdateRow(context.Context, map[string]types.AttributeValue, expression.UpdateBuilder) (*dynamodb.UpdateItemOutput, error)
	DeleteRow(context.Context, map[string]types.AttributeValue) (*dynamodb.DeleteItemOutput, error)
}

type DynamoDB struct {
	*dynamodb.Client
	tableName string
}

func NewDynamoDBClient(ctx context.Context, table string, creds configs.AwsCredentials) (DynamoDbClient, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(creds.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: creds.Key, SecretAccessKey: creds.Secret, SessionToken: "",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return &DynamoDB{
		c,
		table,
	}, nil
}

func (d *DynamoDB) AddRow(ctx context.Context, row any) (*dynamodb.PutItemOutput, error) {
	av, err := attributevalue.MarshalMap(row)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map %w", err)
	}

	res, err := d.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put item %w", err)
	}
	return res, nil
}

func (d *DynamoDB) ReadRow(ctx context.Context, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {
	return d.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key:       key,
	})
}

func (d *DynamoDB) UpdateRow(ctx context.Context, updateKey map[string]types.AttributeValue, updateExp expression.UpdateBuilder) (*dynamodb.UpdateItemOutput, error) {
	expr, err := expression.NewBuilder().WithUpdate(updateExp).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression, %w", err)
	}

	res, err := d.UpdateItem(
		ctx,
		&dynamodb.UpdateItemInput{
			TableName:                 aws.String(d.tableName),
			Key:                       updateKey,
			UpdateExpression:          expr.Update(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ReturnValues:              types.ReturnValueNone,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update item %w", err)
	}
	return res, nil
}

func (d *DynamoDB) DeleteRow(ctx context.Context, deleteKey map[string]types.AttributeValue) (*dynamodb.DeleteItemOutput, error) {
	return d.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key:       deleteKey,
	})
}
