package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// GetUserByUsername retrieves a user by username from the database.
//
// Parameters:
//   - ctx: the context used for the database operation.
//   - username: the username used to find the user.
//
// Returns:
//   - The user matching the provided username.
//   - An error if the user cannot be found or the database operation fails.
func (r *sqlRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	span := newrelic.FromContext(ctx).StartSegment("GetUserByUsername_UserRepository")
	defer span.End()

	user := &model.User{}

	err := r.db.WithContext(ctx).Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return user, nil
}
