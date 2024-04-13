package main

import (
	"context"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/arandich/marketplace-id/internal/model"
	"github.com/arandich/marketplace-id/internal/repository"
	"github.com/arandich/marketplace-id/internal/service"
	grpcTransport "github.com/arandich/marketplace-id/internal/transport/grpc"
	"github.com/arandich/marketplace-id/internal/transport/http"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

func runApp(ctx context.Context, cfg config.Config) {
	logger := zerolog.Ctx(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt, os.Kill)

	// Prometheus.
	promMetrics := initMetrics(ctx, cfg.Prometheus)

	// HTTP
	httpLis, err := initHTTP(ctx, cfg.HTTP)
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to HTTP server")
	}
	defer func() {
		if err = httpLis.Close(); err != nil {
			logger.Error().Err(err).Msg("error closing HTTP listener")
		}
	}()

	srv := http.NewServer(httpLis, cfg.HTTP)
	httpErrCh := srv.StartAndServe()

	// GRPC
	grpcLis, err := initGRPC(ctx, cfg.GRPC)
	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing GRPC listener")
	}
	defer func() {
		if err = grpcLis.Close(); err != nil {
			logger.Error().Err(err).Msg("error closing GRPC listener")
		}
	}()

	pgPool, err := initPostgres(ctx, cfg.Postgres)
	if err != nil {
		logger.Fatal().Err(err).Msg("error initializing postgres")
	}
	defer pgPool.Close()

	//clients := model.Clients{}
	services := model.Services{
		IdService: service.NewIdService(repository.NewIdRepository(ctx, pgPool, promMetrics, cfg)),
	}
	// GRPC.
	grpcTrSrv := grpcTransport.New(ctx, cfg.GRPC)
	grpcSrv, grpcErrCh := grpcTrSrv.Start(ctx, grpcLis, services, promMetrics)
	defer grpcSrv.GracefulStop()

	logger.Info().Str("service", cfg.App.Name).Msg("service started")

	for {
		select {
		case err = <-grpcErrCh:
			logger.Error().Err(err).Msg("retrieved error from GRPC server")
			c <- os.Kill
		case err = <-httpErrCh:
			logger.Error().Err(err).Msg("retrieved error from HTTP server")
			c <- os.Kill
		case sig := <-c:
			logger.Warn().Str("signal", sig.String()).Msg("received shutdown signal")
			return
		}
	}
}
