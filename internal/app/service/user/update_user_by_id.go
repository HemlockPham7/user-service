package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// UpdateUserByID updates the user's profile information by user ID.
//
// Parameters:
//   - ctx: the context used to control the request lifecycle and propagate tracing information.
//   - uid: the unique ID of the user to update.
//   - displayName: the new display name of the user.
//   - email: the new email address of the user.
//
// Returns:
//   - An error if the user cannot be updated.
func (s *service) UpdateUserByID(ctx context.Context, uid, displayName, email string) error {
	span := newrelic.FromContext(ctx).StartSegment("UpdateUserByID_UserService")
	defer span.End()

	updatedUser := &model.User{
		DisplayName: displayName,
		Email:       email,
	}

	return s.repo.UpdateUserByID(ctx, uid, updatedUser)
}
