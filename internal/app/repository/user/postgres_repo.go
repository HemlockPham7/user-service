package user

import "gorm.io/gorm"

type sqlRepository struct {
	db *gorm.DB
}

func NewSqlRepository(db *gorm.DB) Repository {
	return &sqlRepository{db: db}
}
