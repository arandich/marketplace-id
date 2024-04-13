package service

import (
	"context"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
)

type IdRepository interface {
	Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error)
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
