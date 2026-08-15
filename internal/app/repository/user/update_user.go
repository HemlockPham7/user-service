package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
)

func (r *sqlRepository) UpdateUserByID(ctx context.Context, id string, updatedUser *model.User) error {
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updatedUser).Error
	if err != nil {
		return dbutils.CatchDBError(err)
	}

	return nil
}
