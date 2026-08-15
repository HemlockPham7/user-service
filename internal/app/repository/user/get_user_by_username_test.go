package user

import (
	"testing"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputUsername string

		expectedError error
		expectedUser  *model.User
	}{
		{
			name: "normal case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUsername: "johndoe",

			expectedUser: &model.User{
				Base: model.Base{
					ID:        "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
					CreatedAt: fixtures.TestTime,
					UpdatedAt: fixtures.TestTime,
				},
				Username:    "johndoe",
				DisplayName: "John Doe",
				Password:    "johndoe",
				Email:       "johndoe@gmail.com",
			},

			expectedError: nil,
		},
		{
			name: "Get user by username failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUsername: "unknown user",

			expectedUser: nil,

			expectedError: dbutils.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewSqlRepository(db)

			user, err := repo.GetUserByUsername(ctx, tc.inputUsername)

			assert.ErrorIs(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)

		})
	}
}
