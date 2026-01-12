package model

import "time"

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	//Identity fields
	FirstName string
	LastName  string
	Email     string
	Phone     *string

	//Account state
	IsActive        bool
	IsEmailVerified bool
	IsPhoneVerified bool
	IsBanned        bool
	BannedAt        *time.Time
	BanReason       *string

	//Security
	PasswordHash string
	LastLoginAt  *time.Time

	//Authorization
	Role Role
}
