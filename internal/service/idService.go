package service

import (
	"context"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IdRepository interface {
	Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error)
	InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error)
	GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error)
}

var _ IdRepository = (*IdService)(nil)

type IdService struct {
	pb.UnimplementedIdServiceServer
	repository IdRepository
}

func NewIdService(repository IdRepository) IdService {
	return IdService{
		repository: repository,
	}
}

func (s IdService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return s.repository.Auth(ctx, req)
}

func (s IdService) InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error) {
	return s.repository.InitHold(ctx, req)
}

func (s IdService) GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error) {
	return s.repository.GetUser(ctx, req)
}
