package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// GetSelfInfo retrieves the authenticated user's information by user ID.
//
// Parameters:
//   - ctx: the context used to control the request lifecycle and propagate tracing information.
//   - uid: the unique ID of the user whose information should be retrieved.
//
// Returns:
//   - A user model containing the user's information.
//   - An error if the user cannot be retrieved.
func (s *service) GetSelfInfo(ctx context.Context, uid string) (*model.User, error) {
	span := newrelic.FromContext(ctx).StartSegment("GetSelfInfo_UserService")
	defer span.End()

	return s.repo.GetUserByID(ctx, uid)
}
