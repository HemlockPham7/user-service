package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HemlockPham7/user-service/internal/app/model"
	mockUserService "github.com/HemlockPham7/user-service/internal/app/service/user/mocks"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetSelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mockUserService.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful get self info",

			setupRequest: func(ctx *gin.Context) {
				setupGetSelfInfoRequest(ctx, true)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("GetSelfInfo", ctx, "user-123").Return(&model.User{
					Base:        fixtures.GetTestBase("user-123"),
					DisplayName: "user 123",
					Username:    "user123",
					Password:    "password",
					Email:       "user123@example.com",
				}, nil)
				return mockService
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"id":"user-123","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","display_name":"user 123","username":"user123","email":"user123@example.com"}`,
		},
		{
			name: "missing claim",

			setupRequest: func(ctx *gin.Context) {
				setupGetSelfInfoRequest(ctx, false)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				return mockUserService.NewService(t)
			},

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"claim not exist"}`,
		},
		{
			name: "internal error",

			setupRequest: func(ctx *gin.Context) {
				setupGetSelfInfoRequest(ctx, true)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("GetSelfInfo", ctx, "user-123").Return(nil, assert.AnError)
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
			userHandler := NewHandler(serviceMock)
			userHandler.GetSelfInfo(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func setupGetSelfInfoRequest(ctx *gin.Context, haveClaims bool) {
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
	ctx.Request.Header.Set("Content-Type", "application/json")

	if haveClaims {
		ctx.Set("claims", jwt.MapClaims{
			"sub": "user-123",
		})
	}
}
