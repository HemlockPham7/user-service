package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (r *sqlRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	span := newrelic.FromContext(ctx).StartSegment("CreateUser_UserRepository")
	defer span.End()

	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return newUser, err
}
