package dynamodb

import (
	"context"
	"fmt"
	"mr-tasker/configs"
	"mr-tasker/internal/services/user/model"

	"mr-tasker/internal/storage/api"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const tableName = "UsersTable"

func NewUserStorage(ctx context.Context, creds configs.AwsCredentials) (api.UserStorage, error) {
	client, err := NewDynamoDBClient(ctx, tableName, creds)
	if err != err {
		return nil, fmt.Errorf("failed to create dynamodb user storage")
	}
	return &DynamoDbUserStorage{
		client: client,
	}, nil
}

type DynamoDbUserStorage struct {
	client DynamoDbClient
}

func (d *DynamoDbUserStorage) CreateUser(ctx context.Context, u *model.User) (string, error) {
	userID := uuid.NewString()

	user := User{
		UserID: userID,
		Name:   u.Name,
		Age:    u.Age,
	}

	_, err := d.client.AddRow(ctx, user)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (d *DynamoDbUserStorage) ReadUser(ctx context.Context, id string) (*model.User, error) {
	var user User
	userKey, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	searchMap := map[string]types.AttributeValue{
		"UserID": userKey,
	}

	res, err := d.client.ReadRow(ctx, searchMap)
	if err != nil {
		return nil, fmt.Errorf("failed to read user row %w", err)
	}

	err = attributevalue.UnmarshalMap(res.Item, &user)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:   user.UserID,
		Name: user.Name,
		Age:  user.Age,
	}, nil
}

func (d *DynamoDbUserStorage) UpdateUser(ctx context.Context, u *model.User) error {
	userKey, err := attributevalue.Marshal(u.ID)
	if err != nil {
		return err
	}

	updateKey := map[string]types.AttributeValue{
		"UserID": userKey,
	}

	// update fields
	updateExp := expression.Set(
		expression.Name("Name"),
		expression.Value(u.Name),
	)
	updateExp.Set(
		expression.Name("Age"),
		expression.Value(u.Age),
	)

	_, err = d.client.UpdateRow(ctx, updateKey, updateExp)
	if err != nil {
		return fmt.Errorf("failed to update row %w", err)
	}

	return nil
}

func (d *DynamoDbUserStorage) DeleteUser(ctx context.Context, id string) error {
	userKey, err := attributevalue.Marshal(id)
	if err != nil {
		return fmt.Errorf("failed to marshall key, %w", err)
	}

	deleteKey := map[string]types.AttributeValue{
		"UserID": userKey,
	}

	_, err = d.client.DeleteRow(ctx, deleteKey)
	if err != nil {
		return fmt.Errorf("failed to delete row, %w", err)
	}

	return nil
}
