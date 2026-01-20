package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/razedwell/go-hand/internal/model"
	"github.com/razedwell/go-hand/internal/repository/user"
)

type UserRepo struct {
	db *sql.DB
}

var _ user.Repository = (*UserRepo)(nil)

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT
			id, created_at, updated_at,
			first_name, last_name, email, phone,
			is_active, is_email_verified, is_phone_verified,
			is_banned, banned_at, ban_reason,
			password_hash, last_login_at, role
		FROM users
		WHERE email = $1
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.FirstName, &user.LastName, &user.Email, &user.Phone,
		&user.IsActive, &user.IsEmailVerified, &user.IsPhoneVerified,
		&user.IsBanned, &user.BannedAt, &user.BanReason,
		&user.PasswordHash, &user.LastLoginAt, &user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to query user by email")
	}

	return user, nil
}

func (r *UserRepo) FindUserById(ctx context.Context, id int64) (*model.User, error) {
	const query = `
		SELECT
			id, created_at, updated_at,
			first_name, last_name, email, phone,
			is_active, is_email_verified, is_phone_verified,
			is_banned, banned_at, ban_reason,
			password_hash, last_login_at, role
		FROM users
		WHERE id = $1
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.FirstName, &user.LastName, &user.Email, &user.Phone,
		&user.IsActive, &user.IsEmailVerified, &user.IsPhoneVerified,
		&user.IsBanned, &user.BannedAt, &user.BanReason,
		&user.PasswordHash, &user.LastLoginAt, &user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to query user by id")
	}

	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (
			first_name, last_name, email, phone,
			is_active, is_email_verified, is_phone_verified,
			is_banned, banned_at, ban_reason,
			password_hash, last_login_at, role
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.FirstName, user.LastName, user.Email, user.Phone,
		user.IsActive, user.IsEmailVerified, user.IsPhoneVerified,
		user.IsBanned, user.BannedAt, user.BanReason,
		user.PasswordHash, user.LastLoginAt, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
