package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/arandich/marketplace-id/internal/model"
	"github.com/arandich/marketplace-id/pkg/metrics"
	pb "github.com/arandich/marketplace-proto/api/proto/services"
	"github.com/arandich/marketplace-proto/api/proto/types"
	"github.com/arandich/marketplace-sdk/authorization/hash"
	sdkJwt "github.com/arandich/marketplace-sdk/authorization/jwt"
	scripts "github.com/arandich/marketplace-sdk/database-scripts/generated/postgres/marketplace-id"
	sdkModel "github.com/arandich/marketplace-sdk/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type IdRepository struct {
	pgPool      *pgxpool.Pool
	promMetrics metrics.Metrics
	logger      *zerolog.Logger
	cfg         config.Config
	clients     model.Clients
	queries     *scripts.Queries
}

func NewIdRepository(ctx context.Context, pgPool *pgxpool.Pool, promMetrics metrics.Metrics, cfg config.Config, clients model.Clients) *IdRepository {
	return &IdRepository{
		pgPool:      pgPool,
		promMetrics: promMetrics,
		logger:      zerolog.Ctx(ctx),
		cfg:         cfg,
		clients:     clients,
		queries:     scripts.New(pgPool),
	}
}

func (i IdRepository) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return &pb.AuthResponse{}, errors.New("username or password is empty")
	}

	user, err := i.queries.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !hash.CheckPasswordHash(req.GetPassword(), user.Password) {
		return &pb.AuthResponse{}, errors.New("wrong password")
	}

	claims := sdkModel.Claims{
		UserID: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    i.cfg.JWT.SignIssuer,
		},
	}

	token, claims, err := sdkJwt.CreateJWTToken(ctx, claims)
	if err != nil {
		return &pb.AuthResponse{}, fmt.Errorf("failed to create jwt token: %w", err)
	}

	return &pb.AuthResponse{
		JwtToken:  token,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

func (i IdRepository) InitHold(ctx context.Context, req *pb.InitHoldRequest) (*pb.InitHoldResponse, error) {
	return nil, nil
}

func (i IdRepository) GetUser(ctx context.Context, req *emptypb.Empty) (*pb.GetUserResponse, error) {
	jwtToken := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(jwtToken) == 0 {
		return &pb.GetUserResponse{}, errors.New("jwt token is empty")
	}

	claims, err := sdkJwt.ValidateJWTToken(ctx, jwtToken[0])
	if err != nil {
		return &pb.GetUserResponse{}, err
	}

	user, err := i.queries.GetUser(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &pb.GetUserResponse{User: &types.User{
		UserId:  user.UserID,
		Balance: user.Balance,
	}}, nil
}

func (i IdRepository) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	// validate
	if req.GetUsername() == "" {
		return nil, errors.New("username is empty")
	}
	if req.GetPassword() == "" {
		return nil, errors.New("password is empty")
	}
	var balance int64 = 0
	if req.Balance != nil {
		balance = *req.Balance
	}
	userID := uuid.New()

	passwdHash, err := hash.HashPassword(req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	err = i.queries.CreateUser(ctx, scripts.CreateUserParams{
		UserID:   userID.String(),
		Balance:  balance,
		Username: req.GetUsername(),
		Password: passwdHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	claims := sdkModel.Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    i.cfg.JWT.Issuer,
		},
	}
	_, _, err = sdkJwt.CreateJWTToken(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt token: %w", err)
	}

	return &pb.RegisterUserResponse{
		User: &types.User{
			UserId:  userID.String(),
			Balance: balance,
		},
	}, nil
}
