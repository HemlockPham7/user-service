package user

import (
	"context"
	"errors"
	"time"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const tokenDuration = time.Hour * 24

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrCannotGenerateToken = errors.New("cannot generate token")

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	span := newrelic.FromContext(ctx).StartSegment("Login_UserService")
	defer span.End()

	// check user exists with username
	user, err := s.repo.GetUserByUsername(ctx, username)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFound):
		return "", ErrInvalidCredentials
	case err == nil:
	default:
		return "", err
	}

	// compare password hash
	if !s.hasher.Compare(user.Password, password) {
		return "", ErrInvalidCredentials
	}

	// if match -> generate token
	tokenContent := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(tokenDuration).Unix(),
	}
	tokenString, err := s.jwtGenerator.GenerateJWT(tokenContent)
	if err != nil {
		return "", ErrCannotGenerateToken
	}

	// return token
	return tokenString, nil
}
