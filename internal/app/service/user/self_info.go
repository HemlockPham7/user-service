package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (s *service) GetSelfInfo(ctx context.Context, uid string) (*model.User, error) {
	span := newrelic.FromContext(ctx).StartSegment("GetSelfInfo_UserService")
	defer span.End()

	return s.repo.GetUserByID(ctx, uid)
}
