package model

type User struct {
	Base
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"`
	Email       string `gorm:"unique;column:email" json:"email"`
}
