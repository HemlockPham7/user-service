package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockUserService "github.com/HemlockPham7/user-service/internal/app/service/user/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx *gin.Context) *mockUserService.Service

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful update profile",

			setupRequest: func(ctx *gin.Context) {
				setupUpdateUserRequest(ctx, "updated user", "updateduser@example.com", true)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("UpdateUserByID", ctx, "user-123", "updated user", "updateduser@example.com").
					Return(nil)
				return mockService
			},

			expectedCode:     http.StatusOK,
			expectedResponse: `{"message":"User updated successfully"}`,
		},
		{
			name: "missing claim",

			setupRequest: func(ctx *gin.Context) {
				setupUpdateUserRequest(ctx, "updated user", "updateduser@example.com", false)
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
				setupUpdateUserRequest(ctx, "updated user", "updateduser@example.com", true)
			},

			setupMockSvc: func(ctx *gin.Context) *mockUserService.Service {
				mockService := mockUserService.NewService(t)
				mockService.On("UpdateUserByID", ctx, "user-123", "updated user", "updateduser@example.com").
					Return(assert.AnError)
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
			userHandler.UpdateUserByID(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func setupUpdateUserRequest(ctx *gin.Context, displayName, email string, haveClaims bool) {
	inputUpdateUser := &updateUserRequest{
		DisplayName: displayName,
		Email:       email,
	}
	reqBody, _ := json.Marshal(inputUpdateUser)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/users/update", strings.NewReader(string(reqBody)))
	ctx.Request.Header.Set("Content-Type", "application/json")

	if haveClaims {
		ctx.Set("claims", jwt.MapClaims{
			"sub": "user-123",
		})
	}
}
