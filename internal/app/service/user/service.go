package user

import (
	"context"

	"github.com/HemlockPham7/common-libs/pkg/jwtutils"
	"github.com/HemlockPham7/common-libs/pkg/utils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/HemlockPham7/user-service/internal/app/repository/user"
)

//go:generate mockery --name Service --filename service.go --outpkg mockUserService
type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetSelfInfo(ctx context.Context, uid string) (*model.User, error)
	UpdateUserByID(ctx context.Context, uid, displayName, email string) error
}

type service struct {
	repo         user.Repository
	hasher       utils.Hasher
	jwtGenerator jwtutils.JWTGenerator
}

func NewService(repo user.Repository, hasher utils.Hasher, jwtGen jwtutils.JWTGenerator) Service {
	return &service{repo: repo, hasher: hasher, jwtGenerator: jwtGen}
}
