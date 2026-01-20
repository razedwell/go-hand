package postgres

import (
	"context"
	"database/sql"

	"github.com/razedwell/go-hand/internal/model"
	"github.com/razedwell/go-hand/internal/transport/http/helpers"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (r *TokenRepo) GetRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, revoked_at, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1`
	row := r.db.QueryRowContext(ctx, query, tokenHash)

	var rt model.RefreshToken
	var revokedAt sql.NullTime

	err := row.Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &revokedAt, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}

	if revokedAt.Valid {
		rt.RevokedAt = &revokedAt.Time
	} else {
		rt.RevokedAt = nil
	}

	return &rt, nil
}

func (r *TokenRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	now := helpers.GetCurrentTimeStampUTC()
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2`
	_, err := r.db.ExecContext(ctx, query, now, tokenHash)
	return err
}

func (r *TokenRepo) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	now := helpers.GetCurrentTimeStampUTC()
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, now, userID)
	return err
}

func (r *TokenRepo) DeleteExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW() - INTERVAL '1 day'`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
