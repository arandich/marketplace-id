package main

import (
	"context"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/arandich/marketplace-id/pkg/metrics"
	sdkPg "github.com/arandich/marketplace-sdk/postgres"
	sdkPrometheus "github.com/arandich/marketplace-sdk/prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"net"
)

func initHTTP(ctx context.Context, cfg config.HttpConfig) (net.Listener, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing HTTP server listener")

	lis, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func initGRPC(ctx context.Context, cfg config.GrpcConfig) (net.Listener, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing GRPC server listener")

	lis, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func initMetrics(ctx context.Context, cfg config.PrometheusConfig) metrics.Metrics {
	logger := zerolog.Ctx(ctx)
	logger.Info().Str("namespace", cfg.Namespace).Str("subsystem", cfg.Subsystem).Msg("initializing prometheus metrics")

	promCfg := sdkPrometheus.Config{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
	}

	baseMetrics := sdkPrometheus.New(promCfg)
	promMetrics := metrics.New(baseMetrics, cfg)

	return promMetrics
}

func initPostgres(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("initializing postgres connection")

	pgCfg := sdkPg.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
	}

	pool, err := sdkPg.Connect(ctx, pgCfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
