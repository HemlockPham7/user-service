package fixtures

import (
	"github.com/HemlockPham7/user-service/internal/app/model"
	"gorm.io/gorm"
)

type UserCommonTestDB struct {
	base // composition in golang
}

func (u *UserCommonTestDB) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserCommonTestDB) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []*model.User{
		{
			Base:        GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd4"),
			DisplayName: "John Doe",
			Username:    "johndoe",
			Password:    "johndoe",
			Email:       "johndoe@gmail.com",
		},
		{
			Base:        GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd5"),
			DisplayName: "Jane Doe",
			Username:    "janedoe",
			Password:    "janedoe",
			Email:       "janedoe@gmail.com",
		},
		{
			Base:        GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7b10"),
			DisplayName: "Test User",
			Username:    "testuser001",
			Password:    "$2a$10$hhuB9rZrp5ikmRb5yAF9hev6AE2tC404jhtP.bdOjme9lECJClzFu",
			Email:       "test@gmail.com",
		},
	}

	return db.CreateInBatches(users, 10).Error
}
