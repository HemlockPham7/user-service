package infrastructure

import (
	"github.com/HemlockPham7/common-libs/pkg/logger"
	"github.com/HemlockPham7/user-service/internal/api"
	"github.com/gin-gonic/gin"
)

// CreateAPIConfig creates the application configuration from environment variables.
//
// It panics if the configuration cannot be initialized.
func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

// CreateAPI initializes the application and its required dependencies.
//
// It configures the log level, Redis client, database connection, JWT providers,
// New Relic client, and Gin HTTP engine before creating the application engine.
//
// Returns:
//   - The initialized application engine.
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

	// Create New Relic client
	nrClient := CreateNRClient("user")

	app := gin.Default()

	return api.NewEngine(&api.EngineOpts{
		App:         app,
		Cfg:         cfg,
		RedisClient: redisClient,
		DbClient:    db,
		JwtGen:      jwtGen,
		JwtVal:      jwtVal,
		NrClient:    nrClient,
	})
}
