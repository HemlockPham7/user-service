package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (s *service) UpdateUserByID(ctx context.Context, uid, displayName, email string) error {
	span := newrelic.FromContext(ctx).StartSegment("UpdateUserByID_UserService")
	defer span.End()

	updatedUser := &model.User{
		DisplayName: displayName,
		Email:       email,
	}

	return s.repo.UpdateUserByID(ctx, uid, updatedUser)
}
