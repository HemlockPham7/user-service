package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// GetUserByID retrieves a user from the database by their unique ID.
//
// Parameters:
//   - ctx: the context used to control the database operation and propagate tracing information.
//   - id: the unique ID of the user to retrieve.
//
// Returns:
//   - A user model containing the retrieved user's information.
//   - An error if the user does not exist or the database operation fails.
func (r *sqlRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	span := newrelic.FromContext(ctx).StartSegment("GetUserByID_UserRepository")
	defer span.End()

	user := &model.User{}

	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return user, nil
}
