package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
)

func (s *service) UpdateUserByID(ctx context.Context, uid, displayName, email string) error {
	updatedUser := &model.User{
		DisplayName: displayName,
		Email:       email,
	}

	return s.repo.UpdateUserByID(ctx, uid, updatedUser)
}
