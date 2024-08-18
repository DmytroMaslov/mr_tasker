package dynamodb

import (
	"context"
	"testing"

	"mr-tasker/internal/services/user/model"
	dynamodbMock "mr-tasker/internal/storage/dynamodb/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := dynamodbMock.NewMockDynamoDbClient(ctrl)
	defer ctrl.Finish()

	m.EXPECT().AddRow(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, user User) (string, error) {
		assert.NotEmpty(t, user.UserID)
		assert.NotEmpty(t, user.Name)
		assert.NotEmpty(t, user.Age)
		return user.UserID, nil
	})

	user := &model.User{
		Name: "test",
		Age:  18,
	}

	s := &DynamoDbUserStorage{
		client: m,
	}

	res, err := s.CreateUser(context.Background(), user)
	require.Nil(t, err)
	assert.NotEmpty(t, res)
}
