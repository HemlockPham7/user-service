package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

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
