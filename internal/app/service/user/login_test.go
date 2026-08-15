package user

import (
	"context"
	"testing"
	"time"

	mock_jwtgen "github.com/HemlockPham7/common-libs/pkg/jwtutils/mocks"
	mock_hasher "github.com/HemlockPham7/common-libs/pkg/utils/mocks"
	"github.com/HemlockPham7/user-service/internal/app/model"
	mock_user "github.com/HemlockPham7/user-service/internal/app/repository/user/mocks"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockUserRepository func(ctx context.Context) *mock_user.Repository
		setupMockPasswordHash   func(t *testing.T) *mock_hasher.Hasher
		setupMockJWTGen         func(t *testing.T) *mock_jwtgen.JWTGenerator

		inputUsername string
		inputPassword string

		expectedError  error
		expectedOutput string
	}{
		{
			name: "Login successfully",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				expectedUser := &model.User{
					Base:        fixtures.GetTestBase("d7c13097-a60e-4eae-67a7-0b9b533b7bd5"),
					DisplayName: "abc xyz",
					Username:    "abcxyz",
					Password:    "abcxyz",
					Email:       "abcxyz@example.com",
				}
				repoMock.On("GetUserByUsername", ctx, "abcxyz").Return(expectedUser, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Compare", "abcxyz", "abcxyz").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				jwtGeneratorMock := mock_jwtgen.NewJWTGenerator(t)
				tokenContent := jwt.MapClaims{
					"sub":   "d7c13097-a60e-4eae-67a7-0b9b533b7bd5",
					"email": "abcxyz@example.com",
					"iat":   time.Now().Unix(),
					"exp":   time.Now().Add(tokenDuration).Unix(),
				}
				jwtGeneratorMock.On("GenerateJWT", tokenContent).Return("mocked_jwt_token", nil)
				return jwtGeneratorMock
			},

			inputUsername:  "abcxyz",
			inputPassword:  "abcxyz",
			expectedError:  nil,
			expectedOutput: "mocked_jwt_token",
		},
		{
			name: "Fail to get user by username",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "janedoe").Return(nil, ErrInvalidCredentials)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				return mock_hasher.NewHasher(t)
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				return mock_jwtgen.NewJWTGenerator(t)
			},

			inputUsername:  "janedoe",
			inputPassword:  "janedoe",
			expectedError:  ErrInvalidCredentials,
			expectedOutput: "",
		},
		{
			name:           "invalid password",
			inputUsername:  "1234567",
			inputPassword:  "1234567",
			expectedError:  ErrInvalidCredentials,
			expectedOutput: "",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "1234567").Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "1234567",
					Password:    "1234567",
					Email:       "1234567@gmail.com",
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Compare", "1234567", "1234567").Return(false)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				return mock_jwtgen.NewJWTGenerator(t)
			},
		},
		{
			name:           "Fail to generate JWT",
			expectedError:  ErrCannotGenerateToken,
			expectedOutput: "",
			inputUsername:  "111111111",
			inputPassword:  "111111111",

			setupMockUserRepository: func(ctx context.Context) *mock_user.Repository {
				repoMock := mock_user.NewRepository(t)
				repoMock.On("GetUserByUsername", ctx, "111111111").Return(&model.User{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
					DisplayName: "Jane Doe",
					Username:    "111111111",
					Password:    "111111111",
					Email:       "111111111@gmail.com",
				}, nil)
				return repoMock
			},

			setupMockPasswordHash: func(t *testing.T) *mock_hasher.Hasher {
				hashingMock := mock_hasher.NewHasher(t)
				hashingMock.On("Compare", "111111111", "111111111").Return(true)
				return hashingMock
			},

			setupMockJWTGen: func(t *testing.T) *mock_jwtgen.JWTGenerator {
				jwtGeneratorMock := mock_jwtgen.NewJWTGenerator(t)
				jwtGeneratorMock.On("GenerateJWT", jwt.MapClaims{
					"sub":   "d7c13097-67a7-4eae-a60e-0b9b533b7bd5",
					"email": "111111111@gmail.com",
					"iat":   time.Now().Unix(),
					"exp":   time.Now().Add(tokenDuration).Unix(),
				}).Return("", ErrCannotGenerateToken)
				return jwtGeneratorMock
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			userRepoMock := tc.setupMockUserRepository(ctx)
			passwordHasherMock := tc.setupMockPasswordHash(t)
			jwtGeneratorMock := tc.setupMockJWTGen(t)

			userService := NewService(userRepoMock, passwordHasherMock, jwtGeneratorMock)

			output, err := userService.Login(ctx, tc.inputUsername, tc.inputPassword)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
