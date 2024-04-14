package repository_test

import (
	"context"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/arandich/marketplace-id/internal/model"
	"github.com/arandich/marketplace-id/internal/repository"
	"github.com/arandich/marketplace-id/internal/storage/mocks"
	"github.com/arandich/marketplace-id/pkg/metrics"
	sdkJwt "github.com/arandich/marketplace-sdk/authorization/jwt"
	id "github.com/arandich/marketplace-sdk/database-scripts/generated/postgres/marketplace-id"
	sdkModel "github.com/arandich/marketplace-sdk/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestGetUserSuccess(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	m := mocks.NewMockStorage(c)

	userID := uuid.New().String()
	token, _, err := sdkJwt.CreateJWTTokenWithoutCookie(context.Background(), sdkModel.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "marketplace",
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", token))

	m.EXPECT().GetUser(gomock.Any(), userID).Return(id.User{
		UserID:    userID,
		Balance:   0,
		Username:  "test",
		Password:  "test",
		CreatedAt: pgtype.Timestamp{},
	}, nil)

	r := repository.NewIdRepository(ctx, m, metrics.Metrics{}, config.Config{}, model.Clients{})

	resp, err := r.GetUser(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, userID, resp.User.UserId)
}

func TestGetUserInvalidIssuer(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	userID := uuid.New().String()

	token, _, err := sdkJwt.CreateJWTTokenWithoutCookie(context.Background(), sdkModel.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "invalid",
		},
	})
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestGetUserNotFound(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	m := mocks.NewMockStorage(c)

	userID := uuid.New().String()

	token, _, err := sdkJwt.CreateJWTTokenWithoutCookie(context.Background(), sdkModel.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "marketplace",
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", token))

	m.EXPECT().GetUser(gomock.Any(), userID).Return(id.User{}, pgx.ErrNoRows)

	r := repository.NewIdRepository(ctx, m, metrics.Metrics{}, config.Config{}, model.Clients{})

	resp, err := r.GetUser(ctx, nil)

	assert.ErrorContains(t, err, "user not found")
	assert.Nil(t, resp)
}

func TestGetUserEmptyJwt(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	m := mocks.NewMockStorage(c)

	r := repository.NewIdRepository(context.Background(), m, metrics.Metrics{}, config.Config{}, model.Clients{})

	resp, err := r.GetUser(context.Background(), nil)

	assert.ErrorContains(t, err, "jwt token is empty")
	assert.Nil(t, resp)
}
