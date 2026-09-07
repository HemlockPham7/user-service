package user

import "gorm.io/gorm"

type sqlRepository struct {
	db *gorm.DB
}

// NewSqlRepository creates a user repository backed by the provided GORM database client.
//
// Parameters:
//   - db: the GORM database client used to persist user data.
//
// Returns:
//   - A user repository backed by the provided database client.
func NewSqlRepository(db *gorm.DB) Repository {
	return &sqlRepository{db: db}
}
