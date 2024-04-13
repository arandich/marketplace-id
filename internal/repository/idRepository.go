package repository

import (
	"context"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/arandich/marketplace-id/pkg/metrics"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type IdRepository struct {
	pgPool      *pgxpool.Pool
	promMetrics metrics.Metrics
	logger      *zerolog.Logger
	cfg         config.Config
}

func NewIdRepository(ctx context.Context, pgPool *pgxpool.Pool, promMetrics metrics.Metrics, cfg config.Config) *IdRepository {
	return &IdRepository{
		pgPool:      pgPool,
		promMetrics: promMetrics,
		logger:      zerolog.Ctx(ctx),
		cfg:         cfg,
	}
}

func (i IdRepository) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		JwtToken: "test",
	}, nil
}
