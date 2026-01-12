package model

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string

	ExpiresAt time.Time
	RevokedAt *time.Time

	CreatedAt time.Time
}

type VerificationType string

const (
	VerifyEmail   VerificationType = "email"
	VerifyPhone   VerificationType = "phone"
	PasswordReset VerificationType = "password_reset"
)

type VerificationCode struct {
	ID        int64
	UserID    int64
	CodeHash  string
	Type      VerificationType // email, phone, password_reset
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
