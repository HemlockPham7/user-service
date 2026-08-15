package user

import (
	"testing"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputUserID string

		expectedError error
		expectedUser  *model.User
	}{
		{
			name: "normal case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUserID: "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",

			expectedError: nil,

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
		},
		{
			name: "Get user by ID failed - user not found",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUserID: "d7c13097-67a7-4eae-a60e-0b9b533b7b36",

			expectedError: dbutils.ErrRecordNotFound,

			expectedUser: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewSqlRepository(db)

			user, err := repo.GetUserByID(ctx, tc.inputUserID)

			assert.ErrorIs(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)

		})
	}
}
