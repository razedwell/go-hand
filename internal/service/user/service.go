package user

import (
	"context"

	"github.com/razedwell/go-hand/internal/model"
	"github.com/razedwell/go-hand/internal/repository/user"
	"github.com/razedwell/go-hand/internal/security"
	"github.com/razedwell/go-hand/internal/transport/http/helpers"
)

type RegParams struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
	Password  string `json:"password" validate:"required,min=8"`
}

type Service struct {
	users user.Repository
	// Add fields if necessary
}

func NewService(users user.Repository) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user RegParams) error {
	now := helpers.GetCurrentTimeStampUTC()
	// Hash the password
	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Create a new user model
	newUser := &model.User{
		CreatedAt: now,
		UpdatedAt: now,

		// Identity fields
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Phone:        &user.Phone,
		PasswordHash: hashedPassword,

		// Account state
		IsActive:        true,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		IsBanned:        false,

		// Authorization
		Role: model.RoleUser,
	}

	return s.users.CreateUser(ctx, newUser)
}
