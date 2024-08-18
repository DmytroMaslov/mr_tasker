package api

import (
	"context"
	user "mr-tasker/internal/services/user/model"
)

//go:generate mockgen -source=api.go -destination=mocks/api_mock.go
type UserStorage interface {
	CreateUser(context.Context, *user.User) (string, error)
	ReadUser(context.Context, string) (*user.User, error)
	UpdateUser(context.Context, *user.User) error
	DeleteUser(context.Context, string) error
}
