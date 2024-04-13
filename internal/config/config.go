package config

import (
	"context"
	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	App        AppConfig        `yaml:"app" env:",prefix=APP_"`
	Prometheus PrometheusConfig `yaml:"prometheus" env:",prefix=PROMETHEUS_"`
	HTTP       HttpConfig       `yaml:"HTTP" env:",prefix=HTTP_"`
	GRPC       GrpcConfig       `yaml:"GRPC" env:",prefix=GRPC_"`
	Logger     LoggerConfig     `yaml:"logger" env:",prefix=LOG_"`
	Postgres   PostgresConfig   `yaml:"postgres" env:",prefix=POSTGRES_"`
}

func New(ctx context.Context, cfg interface{}) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = envconfig.Process(ctx, cfg)
	if err != nil {
		return err
	}

	err = defaults.Set(cfg)
	if err != nil {
		return err
	}

	return nil
}
