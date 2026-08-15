package infrastructure

import (
	redisPkg "github.com/HemlockPham7/common-libs/pkg/redis"
	"github.com/HemlockPham7/common-libs/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

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

func CreateRedisClient(envPrefix string) *redis.Client {
	redisClient, err := redisPkg.NewClient(envPrefix)
	if err != nil {
		panic(err)
	}
	return redisClient
}

const migrationPath = "file://./migration"

func MigrateDB(dbClient *gorm.DB) error {
	return sqldb.MigrateSQLDB(dbClient, migrationPath, "up", 0)
}
