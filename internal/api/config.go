package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config struct for configuration
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	LogLevel    string `default:"info" envconfig:"LOG_LEVEL"`
	BasePath    string `default:"/" envconfig:"BASE_PATH"`
	ServiceName string `default:"user-service" envconfig:"SERVICE_NAME"`
	InstanceID  string `default:"" envconfig:"INSTANCE_ID"`
}

// NewConfig creates a new config
func NewConfig() (*Config, error) {
	//err := godotenv.Load()
	//if err != nil {
	//	panic(err)
	//}

	cfg := &Config{}
	err := envconfig.Process("user", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}
	return cfg, err
}
