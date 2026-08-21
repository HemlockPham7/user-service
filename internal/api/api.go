package api

import (
	"fmt"
	"net/http"

	"github.com/HemlockPham7/common-libs/pkg/jwtutils"
	"github.com/HemlockPham7/common-libs/pkg/middleware"
	"github.com/HemlockPham7/common-libs/pkg/ratelimitutils"
	"github.com/HemlockPham7/common-libs/pkg/utils"
	"github.com/HemlockPham7/user-service/docs"
	_ "github.com/HemlockPham7/user-service/docs"
	userHdl "github.com/HemlockPham7/user-service/internal/app/handler/user"
	userRepo "github.com/HemlockPham7/user-service/internal/app/repository/user"
	userSvc "github.com/HemlockPham7/user-service/internal/app/service/user"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Engine interface for starting the application
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct for starting the application
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
	dbClient    *gorm.DB
	jwtGen      jwtutils.JWTGenerator
	jwtVal      jwtutils.JWTValidator
	nrClient    *newrelic.Application
}

type EngineOpts struct {
	App         *gin.Engine
	Cfg         *Config
	RedisClient *redis.Client
	DbClient    *gorm.DB
	JwtGen      jwtutils.JWTGenerator
	JwtVal      jwtutils.JWTValidator
	NrClient    *newrelic.Application
}

// NewEngine creates a new engine
func NewEngine(opts *EngineOpts) Engine {
	app := &engine{
		app:         opts.App,
		cfg:         opts.Cfg,
		redisClient: opts.RedisClient,
		dbClient:    opts.DbClient,
		jwtGen:      opts.JwtGen,
		jwtVal:      opts.JwtVal,
		nrClient:    opts.NrClient,
	}
	app.initRoutes()
	return app
}

// Start starts the application
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP to test the API endpoint
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

type handlers struct {
	userHandler userHdl.Handler
}

func (e *engine) initHandlers() *handlers {
	userRepository := userRepo.NewSqlRepository(e.dbClient)
	hasher := utils.NewHasher()
	userService := userSvc.NewService(userRepository, hasher, e.jwtGen)
	userHandler := userHdl.NewHandler(userService)

	return &handlers{
		userHandler: userHandler,
	}
}

type middlewares struct {
	jwtAuth   middleware.JWTAuth
	rateLimit middleware.RateLimit
}

// initMiddlewares initializes the middlewares
func (e *engine) initMiddlewares() middlewares {
	jwtAuth := middleware.NewJWTAuth(e.jwtVal)

	rateLimitRepository := ratelimitutils.NewRedisRepo(e.redisClient)
	rateLimit := middleware.NewRateLimit(rateLimitRepository)

	return middlewares{
		jwtAuth:   jwtAuth,
		rateLimit: rateLimit,
	}
}

// initRoutes initializes the routes
func (e *engine) initRoutes() {
	allHandlers := e.initHandlers()
	allMiddlewares := e.initMiddlewares()

	// Add New Relic middleware
	e.app.Use(nrgin.Middleware(e.nrClient))

	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	privateRoutes := e.app.Group("")
	privateRoutes.Use(allMiddlewares.jwtAuth.JWTAuth())
	privateRoutes.Use(allMiddlewares.rateLimit.RateLimit())
	{
		privateV1Routes := privateRoutes.Group("/v1")
		{
			selfRoutes := privateV1Routes.Group("/self")
			{
				selfRoutes.GET("/info", allHandlers.userHandler.GetSelfInfo)
			}

			userRoutes := privateV1Routes.Group("/users")
			{
				userRoutes.PUT("/update", allHandlers.userHandler.UpdateUserByID)
			}
		}
	}

	publicRoutes := e.app.Group("")
	{
		publicV1Routes := publicRoutes.Group("/v1")
		{
			usersRoutes := publicV1Routes.Group("/users")
			{
				usersRoutes.POST("/register", allHandlers.userHandler.Register)
				usersRoutes.POST("/login", allHandlers.userHandler.Login)
			}
		}
	}
}
