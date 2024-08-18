package user

import (
	"context"
	"errors"
	"fmt"
	"mr-tasker/internal/services/user/model"
	"mr-tasker/internal/storage/api"
)

var (
	ErrNotFound = errors.New("not found")
)

type UserService interface {
	Read(string) (*model.User, error)
	Create(*model.User) (string, error)
	Update(*model.User) (*model.User, error)
	Delete(string) error
}

type UserServiceImp struct {
	storage api.UserStorage
}

func NewUserService(storage api.UserStorage) UserService {
	return &UserServiceImp{
		storage: storage,
	}
}

func (u *UserServiceImp) Read(id string) (*model.User, error) {
	user, err := u.storage.ReadUser(context.TODO(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to read user, %w", err)
	}
	if user.ID == "" {
		return nil, ErrNotFound
	}
	return user, nil

}

func (u *UserServiceImp) Create(user *model.User) (string, error) {
	id, err := u.storage.CreateUser(context.TODO(), user)
	if err != nil {
		return "", fmt.Errorf("failed to create user")
	}
	return id, nil
}

func (u *UserServiceImp) Update(user *model.User) (*model.User, error) {
	err := u.storage.UpdateUser(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user %w", err)
	}

	return u.storage.ReadUser(context.TODO(), user.ID)
}

func (u *UserServiceImp) Delete(id string) error {
	return u.storage.DeleteUser(context.TODO(), id)
}
