package user

import (
	"context"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// CreateUser creates a new user after hashing the provided password.
//
// Parameters:
//   - ctx: the context used for the operation.
//   - username: the username of the new user.
//   - password: the plain-text password to hash before storing.
//   - displayName: the display name of the new user.
//   - email: the email address of the new user.
//
// Returns:
//   - The created user, or an error if password hashing or user creation fails.
func (s *service) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	span := newrelic.FromContext(ctx).StartSegment("CreateUser_UserService")
	defer span.End()

	// hash password
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// create user model
	newUser := &model.User{
		Username:    username,
		Password:    hash,
		Email:       email,
		DisplayName: displayName,
	}

	// call repo to create user
	res, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// return user
	return res, nil
}
