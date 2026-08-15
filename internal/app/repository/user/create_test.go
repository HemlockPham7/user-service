package user

import (
	"testing"

	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/HemlockPham7/user-service/internal/integration_test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputUserName string
		inputEmail    string

		expectedError error
		verifyFunc    func(db *gorm.DB, userName, email string)
	}{
		{
			name: "normal case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},

			inputUserName: "test1",
			inputEmail:    "test1@gmail.com",

			expectedError: nil,
			verifyFunc: func(db *gorm.DB, userName, email string) {
				user := model.User{}
				err := db.First(&user, "username = ?", userName).Error
				assert.NoError(t, err)
				assert.Equal(t, userName, user.Username)
				assert.Equal(t, email, user.Email)
				assert.NotEmpty(t, user.ID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewSqlRepository(db)

			user, err := repo.CreateUser(ctx, &model.User{
				Username: tc.inputUserName,
				Email:    tc.inputEmail,
			})
			if err != nil {
				assert.NotNil(t, user)
			}
			assert.ErrorIs(t, tc.expectedError, err)
			tc.verifyFunc(db, tc.inputUserName, tc.inputEmail)
		})
	}
}
