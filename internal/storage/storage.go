package storage

import (
	"context"
	scripts "github.com/arandich/marketplace-sdk/database-scripts/generated/postgres/marketplace-id"
)

//go:generate mockgen -destination=./mocks/mockStorage.go -package=mocks github.com/arandich/marketplace-id/internal/storage Storage
type Storage interface {
	CreateHold(ctx context.Context, arg scripts.CreateHoldParams) error
	CreateUser(ctx context.Context, arg scripts.CreateUserParams) error
	GetHold(ctx context.Context, holdID int64) (scripts.GetHoldRow, error)
	GetHolds(ctx context.Context, userID string) ([]scripts.GetHoldsRow, error)
	GetUser(ctx context.Context, userID string) (scripts.User, error)
	UpdateHoldStatus(ctx context.Context, arg scripts.UpdateHoldStatusParams) error
	UpdateUserBalance(ctx context.Context, arg scripts.UpdateUserBalanceParams) error
}
