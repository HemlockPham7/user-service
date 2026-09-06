package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// UpdateUserByID updates a user's information in the database by user ID.
//
// Parameters:
//   - ctx: the context used to control the database operation and propagate tracing information.
//   - id: the unique ID of the user to update.
//   - updatedUser: the user data containing the fields to be updated.
//
// Returns:
//   - An error if the database operation fails or the user cannot be found.
//   - Nil if the user is successfully updated.
func (r *sqlRepository) UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error {
	span := newrelic.FromContext(ctx).StartSegment("UpdateUserByID_UserRepository")
	defer span.End()

	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updatedUser)

	if result.Error != nil {
		return dbutils.CatchDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return dbutils.ErrRecordNotFoundType
	}

	return nil
}
