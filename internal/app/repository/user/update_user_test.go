package user

import (
	"testing"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputID       string
		inputUserData *model.User

		expectedError error
	}{
		{
			name: "normal case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputID: "d7c13097-67a7-4eae-a60e-0b9b533b7b10",
			inputUserData: &model.User{
				DisplayName: "User updated",
				Email:       "userupdated@gmail.com",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewSqlRepository(db)

			err := repo.UpdateUserByID(ctx, tc.inputID, tc.inputUserData)

			assert.ErrorIs(t, tc.expectedError, err)
		})
	}
}
