package user

import (
	"context"
	"testing"

	mock_hasher "github.com/HemlockPham7/common-libs/pkg/utils/mocks"
	"github.com/HemlockPham7/user-service/internal/app/model"
	mock_user "github.com/HemlockPham7/user-service/internal/app/repository/user/mocks"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepository func(ctx context.Context) *mock_user.Repository
		setupMockHasher         func(t *testing.T) *mock_hasher.Hasher

		inputUsername    string
		inputPassword    string
		inputDisplayName string
		inputEmail       string

		expectedError  error
		expectedOutput *model.User
	}{
		{
			name: "Create user successfully",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				newUser := &model.User{
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}
				repoMock.On("CreateUser", ctx, newUser).Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}, nil)
				return repoMock
			},

			setupMockHasher: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Hash", "janedoe").Return("janedoe", nil)
				return hashingMock
			},

			inputUsername:    "janedoe",
			inputPassword:    "janedoe",
			inputDisplayName: "Jane Doe",
			inputEmail:       "janedoe@gmail.com",

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
			name: "failed due to create user error",

			setupMockHasher: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Hash", "janedoe").Return("janedoe", nil)
				return hashingMock
			},

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				newUser := &model.User{
					DisplayName: "Jane Doe",
					Username:    "janedoe",
					Password:    "janedoe",
					Email:       "janedoe@gmail.com",
				}
				repoMock.On("CreateUser", ctx, newUser).Return(nil, assert.AnError)
				return repoMock
			},

			inputUsername:    "janedoe",
			inputPassword:    "janedoe",
			inputDisplayName: "Jane Doe",
			inputEmail:       "janedoe@gmail.com",

			expectedError:  assert.AnError,
			expectedOutput: nil,
		},
		{
			name: "fail due to hash error",

			setupMockHasher: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Hash", "janedoe").Return("", assert.AnError)
				return hashingMock
			},

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				return mock_user.NewRepository(t)
			},

			inputUsername:    "janedoe",
			inputPassword:    "janedoe",
			inputDisplayName: "Jane Doe",
			inputEmail:       "janedoe@gmail.com",

			expectedError:  assert.AnError,
			expectedOutput: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			userRepoMock := tc.setupMockUserRepository(ctx)
			hasherMock := tc.setupMockHasher(t)

			userService := NewService(userRepoMock, hasherMock, nil)

			output, err := userService.CreateUser(ctx, tc.inputUsername, tc.inputPassword, tc.inputDisplayName, tc.inputEmail)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
