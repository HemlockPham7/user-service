package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
)

func (s *service) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	// hash password
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// create user model
	newUser := &model.User{
		Username:    username,
		Password:    hash,
		Email:       email,
		DisplayName: displayName,
	}

	// call repo to create user
	res, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// return user
	return res, nil
}
