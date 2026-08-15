package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
)

//go:generate mockery --name Repository --filename user_repo.go --outpkg mock_user
type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error
}
