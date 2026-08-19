package infrastructure

import (
	"github.com/HemlockPham7/common-libs/pkg/logger"
	"github.com/HemlockPham7/user-service/internal/api"
	"github.com/gin-gonic/gin"
)

func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func CreateAPI() api.Engine {
	// create app config
	cfg := CreateAPIConfig()

	// set log level
	logger.SetLogLevel(cfg.LogLevel)

	// create redis client
	redisClient := CreateRedisClient("user")

	// Init db
	db := CreateDB("user")

	jwtGen, jwtVal := CreateJWTProvider()

	app := gin.Default()

	return api.NewEngine(&api.EngineOpts{
		App:         app,
		Cfg:         cfg,
		RedisClient: redisClient,
		DbClient:    db,
		JwtGen:      jwtGen,
		JwtVal:      jwtVal,
	})
}
