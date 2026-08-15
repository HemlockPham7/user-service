package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
)

func (s *service) GetSelfInfo(ctx context.Context, uid string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, uid)
}
