package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	mockUserService "github.com/HemlockPham7/user-service/internal/app/service/user/mocks"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var mockRegisterInput = registerInput{
	Username:    "user-123",
	Password:    "password",
	DisplayName: "User 123",
	Email:       "user-123@example.com",
}

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mockUserService.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful register",

			setupRequest: func(ctx *gin.Context) {
				setupRegisterRequest(ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("CreateUser", ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email).
					Return(&model.User{
						Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd4"),
						DisplayName: mockRegisterInput.DisplayName,
						Username:    mockRegisterInput.Username,
						Password:    mockRegisterInput.Password,
						Email:       mockRegisterInput.Email,
					}, nil)
				return mockService
			},

			expectedCode:     http.StatusCreated,
			expectedResponse: `{"data":{"id":"d7c13097-67a7-4eae-a60e-0b9b533b7bd4","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","display_name":"User 123","username":"user-123","email":"user-123@example.com"},"message":"Register an user successfully!"}`,
		},
		{
			name: "invalid input",

			setupRequest: func(ctx *gin.Context) {
				setupRegisterRequest(ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, "")
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				return mockUserService.NewService(t)
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Input error","details":["Email is invalid (required)"]}`,
		},
		{
			name: "username already exists",

			setupRequest: func(ctx *gin.Context) {
				setupRegisterRequest(ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("CreateUser", ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email).Return(nil, dbutils.ErrDuplicationUsername)
				return mockService
			},

			expectedCode:     http.StatusConflict,
			expectedResponse: `{"message":"Username already taken"}`,
		},
		{
			name: "Email already exists",

			setupRequest: func(ctx *gin.Context) {
				setupRegisterRequest(ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("CreateUser", ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email).Return(nil, dbutils.ErrDuplicationEmail)
				return mockService
			},

			expectedCode:     http.StatusConflict,
			expectedResponse: `{"message":"Email already taken"}`,
		},
		{
			name: "Internal server error",

			setupRequest: func(ctx *gin.Context) {
				setupRegisterRequest(ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("CreateUser", ctx, mockRegisterInput.Username, mockRegisterInput.Password, mockRegisterInput.DisplayName, mockRegisterInput.Email).Return(nil, assert.AnError)
				return mockService
			},

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Processing Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)
			mockService := tc.setupMockSvc(ctx)
			userHandler := NewHandler(mockService)
			userHandler.Register(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func setupRegisterRequest(ctx *gin.Context, username, password, displayName, email string) {
	reqBody, _ := json.Marshal(&registerInput{
		Username:    username,
		Password:    password,
		DisplayName: displayName,
		Email:       email,
	})
	ctx.Request = httptest.NewRequest(
		"POST",
		"/v1/users/register",
		strings.NewReader(string(reqBody)),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")
}
