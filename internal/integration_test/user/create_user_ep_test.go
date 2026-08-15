package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	redisPkg "github.com/HemlockPham7/common-libs/pkg/redis"
	"github.com/HemlockPham7/user-service/internal/api"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserEndpoint_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "user register successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequestUserRegister("user1234", "user1234", "User 1234", "user1234@gmail.com")
				api.ServeHTTP(rec, req)
				return rec
			},

			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `message":"Register an user successfully!"`,
		},
		{
			name: "invalid user register payload",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req, rec := setupRequestUserRegister("password_too_short", "short", "User 1234", "user1234@gmail.com")
				api.ServeHTTP(rec, req)
				return rec
			},

			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `"Password is invalid (gte)"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			setupRedisClient := redisPkg.InitMockRedis(t)
			setupDB := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			testAPI := api.NewEngine(&api.EngineOpts{
				App:         gin.Default(),
				Cfg:         &api.Config{},
				RedisClient: setupRedisClient,
				DbClient:    setupDB,
			})
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}
}

func setupRequestUserRegister(username, password, displayName, email string) (*http.Request, *httptest.ResponseRecorder) {
	reqBody := fmt.Sprintf(`{"username":"%s","password":"%s","display_name":"%s","email":"%s"}`, username, password, displayName, email)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	return req, rec
}
