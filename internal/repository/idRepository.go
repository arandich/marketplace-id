package repository

import (
	"context"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/arandich/marketplace-id/internal/model"
	"github.com/arandich/marketplace-id/pkg/metrics"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IdRepository struct {
	pgPool      *pgxpool.Pool
	promMetrics metrics.Metrics
	logger      *zerolog.Logger
	cfg         config.Config
	clients     model.Clients
}

func NewIdRepository(ctx context.Context, pgPool *pgxpool.Pool, promMetrics metrics.Metrics, cfg config.Config, clients model.Clients) *IdRepository {
	return &IdRepository{
		pgPool:      pgPool,
		promMetrics: promMetrics,
		logger:      zerolog.Ctx(ctx),
		cfg:         cfg,
		clients:     clients,
	}
}

func (i IdRepository) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		JwtToken: "test",
	}, nil
}

func (i IdRepository) InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error) {
	return nil, nil
}

func (i IdRepository) GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{}, nil
}

func (i IdRepository) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return &pb.RegisterUserResponse{}, nil
}
