package user

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HemlockPham7/common-libs/pkg/jwtutils/mocks"
	"github.com/HemlockPham7/common-libs/pkg/middleware"
	redisPkg "github.com/HemlockPham7/common-libs/pkg/redis"
	"github.com/HemlockPham7/user-service/internal/api"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestUserEndpoint_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP         func(api api.Engine) *httptest.ResponseRecorder
		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator
		setupRateLimit        func(ctx context.Context, redisClient *redis.Client) *redis.Client

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "user self info successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequestUpdateUser(true)
				api.ServeHTTP(rec, req)
				return rec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				mockJWTValidator := mocks.NewJWTValidator(t)
				mockJWTValidator.On("ValidateJWT", "valid_jwt_token").
					Return(jwt.MapClaims{"sub": "d7c13097-67a7-4eae-a60e-0b9b533b7bd4"}, nil)
				return mockJWTValidator
			},

			setupRateLimit: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return setupRateLimitUpdateUser(ctx, 1, redisClient)
			},

			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"User updated successfully"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setupRedisClient := redisPkg.InitMockRedis(t)
			redisClient := tc.setupRateLimit(context.Background(), setupRedisClient)
			setupDB := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			setupJWTValidator := tc.setupMockJWTValidator(t)
			testAPI := api.NewEngine(&api.EngineOpts{
				App:         gin.Default(),
				Cfg:         &api.Config{},
				RedisClient: redisClient,
				DbClient:    setupDB,
				JwtVal:      setupJWTValidator,
			})
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func setupRequestUpdateUser(haveClaims bool) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPut, "/v1/users/update", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_jwt_token")

	if haveClaims {
		req.Header.Set("claims", "d7c13097-67a7-4eae-a60e-0b9b533b7bd4")
	}

	rec := httptest.NewRecorder()
	return req, rec
}

func setupRateLimitUpdateUser(ctx context.Context, rateLimit int, redisClient *redis.Client) *redis.Client {
	key := fmt.Sprintf(middleware.RateLimitKeyFormat, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4")
	redisClient.Set(ctx, key, rateLimit, middleware.RateLimitInterval)
	return redisClient
}
