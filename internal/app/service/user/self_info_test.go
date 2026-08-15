package user

import (
	"context"
	"testing"

	"github.com/HemlockPham7/user-service/internal/app/model"
	mock_user "github.com/HemlockPham7/user-service/internal/app/repository/user/mocks"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestService_GetSelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepository func(ctx context.Context) *mock_user.Repository

		inputUID string

		expectedError  error
		expectedOutput *model.User
	}{
		{
			name: "Get self info successfully",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByID", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd5").Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}, nil)
				return repoMock
			},

			inputUID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd5",
			expectedError: nil,
			expectedOutput: &model.User{
				Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
				DisplayName: "Jane Doe",
				Username:    "janedoe",
				Password:    "janedoe",
				Email:       "janedoe@gmail.com",
			},
		},
		{
			name: "Get self info failed",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByID", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd5").Return(nil, assert.AnError)
				return repoMock
			},

			inputUID:       "d7c13097-67a7-4eae-a60e-0b9b533b7bd5",
			expectedError:  assert.AnError,
			expectedOutput: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repoMock := tc.setupMockUserRepository(ctx)
			userService := NewService(repoMock, nil, nil)

			userInformation, err := userService.GetSelfInfo(ctx, tc.inputUID)

			assert.ErrorIs(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, userInformation)
		})
	}
}
