package infrastructure

import (
	redisPkg "github.com/HemlockPham7/common-libs/pkg/redis"
	"github.com/HemlockPham7/common-libs/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CreateDB creates a GORM database client and applies the database migrations.
//
// Parameters:
//   - envPrefix: the environment variable prefix used to load database configuration.
//
// Returns:
//   - A configured and migrated GORM database client.
//
// Panics:
//   - If the database client cannot be created or the migrations fail.
func CreateDB(envPrefix string) *gorm.DB {
	dbClient, err := sqldb.NewClient(envPrefix)
	if err != nil {
		panic(err)
	}
	//err = db.AutoMigrate(&model.User{})
	err = MigrateDB(dbClient)
	if err != nil {
		panic(err)
	}
	return dbClient
}

// CreateRedisClient creates a Redis client using configuration loaded from environment variables.
//
// Parameters:
//   - envPrefix: the environment variable prefix used to load Redis configuration.
//
// Returns:
//   - A configured Redis client.
//
// Panics:
//   - If the Redis client cannot be created.
func CreateRedisClient(envPrefix string) *redis.Client {
	redisClient, err := redisPkg.NewClient(envPrefix)
	if err != nil {
		panic(err)
	}
	return redisClient
}

const migrationPath = "file://./migration"

// MigrateDB applies all pending database migrations.
//
// Parameters:
//   - dbClient: the GORM database client used to execute the migrations.
//
// Returns:
//   - An error if the database migrations fail.
func MigrateDB(dbClient *gorm.DB) error {
	return sqldb.MigrateSQLDB(dbClient, migrationPath, "up", 0)
}
