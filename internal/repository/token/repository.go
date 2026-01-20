package token

import (
	"context"

	"github.com/razedwell/go-hand/internal/model"
)

type Repository interface {
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context) error
}
