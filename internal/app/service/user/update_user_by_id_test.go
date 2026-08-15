package user

import (
	"context"
	"testing"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	mock_user "github.com/HemlockPham7/user-service/internal/app/repository/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_UpdateUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepository func(ctx context.Context) *mock_user.Repository

		inputID          string
		inputDisplayName string
		inputEmail       string

		expectedError error
	}{
		{
			name: "Update user by ID successfully",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("UpdateUserByID", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7b10", &model.User{
					DisplayName: "Updated user",
					Email:       "updateduser@example.com",
				}).Return(nil)
				return repoMock
			},

			inputID:          "d7c13097-67a7-4eae-a60e-0b9b533b7b10",
			inputDisplayName: "Updated user",
			inputEmail:       "updateduser@example.com",

			expectedError: nil,
		},
		{
			name: "Fail to update user by ID - user not found",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("UpdateUserByID", ctx, "nonexistedid", &model.User{
					DisplayName: "Updated user",
					Email:       "updateduser@example.com",
				}).Return(dbutils.ErrRecordNotFoundType)
				return repoMock
			},

			inputID:          "nonexistedid",
			inputDisplayName: "Updated user",
			inputEmail:       "updateduser@example.com",

			expectedError: dbutils.ErrRecordNotFoundType,
		},
		{
			name: "Fail to update user by ID - duplicate email",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("UpdateUserByID", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7b10", &model.User{
					DisplayName: "Updated user",
					Email:       "duplicateemail@example.com",
				}).Return(dbutils.ErrDuplicationType)
				return repoMock
			},

			inputID:          "d7c13097-67a7-4eae-a60e-0b9b533b7b10",
			inputDisplayName: "Updated user",
			inputEmail:       "duplicateemail@example.com",

			expectedError: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repoMock := tc.setupMockUserRepository(ctx)
			userService := NewService(repoMock, nil, nil)
			err := userService.UpdateUserByID(ctx, tc.inputID, tc.inputDisplayName, tc.inputEmail)

			assert.ErrorIs(t, tc.expectedError, err)
		})
	}
}
