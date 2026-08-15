package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HemlockPham7/user-service/internal/app/service/user"
	mockUserService "github.com/HemlockPham7/user-service/internal/app/service/user/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupLoginRequest(ctx *gin.Context, username, password string) {
	reqBody, _ := json.Marshal(&loginInput{
		Username: username,
		Password: password,
	})
	ctx.Request = httptest.NewRequest(
		"POST",
		"/v1/users/login",
		strings.NewReader(string(reqBody)),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")
}

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mockUserService.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful login",

			setupRequest: func(ctx *gin.Context) {
				setupLoginRequest(ctx, "user-123", "password")
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("Login", ctx, "user-123", "password").
					Return("mocked-jwt-gen", nil)
				return mockService
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"data":"mocked-jwt-gen","message":"Logged in successfully!"}`,
		},
		{
			name: "invalid input",

			setupRequest: func(ctx *gin.Context) {
				setupLoginRequest(ctx, "", "")
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				return mockUserService.NewService(t)
			},

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Input error","details":["Username is invalid (required)","Password is invalid (required)"]}`,
		},
		{
			name: "user invalid credentials",

			setupRequest: func(ctx *gin.Context) {
				setupLoginRequest(ctx, "user-123", "password")
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("Login", ctx, "user-123", "password").
					Return("", user.ErrInvalidCredentials)
				return mockService
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"invalid credentials"}`,
		},
		{
			name: "internal error server",

			setupRequest: func(ctx *gin.Context) {
				setupLoginRequest(ctx, "user-123", "password")
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("Login", ctx, "user-123", "password").
					Return("", assert.AnError)
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
			serviceMock := tc.setupMockSvc(ctx)

			handler := NewHandler(serviceMock)
			handler.Login(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
