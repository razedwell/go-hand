package user

import (
	"context"
	"errors"

	"github.com/razedwell/go-hand/internal/repository/user"
	"github.com/razedwell/go-hand/internal/security"
)

var errUnauthorized = errors.New("unauthorized")

type Service struct {
	users user.Repository
	// Add fields if necessary
}

func NewService(users user.Repository) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.users.FindUserByEmail(email)
	if err != nil || security.VerifyPassword(user.PasswordHash, password) == false {
		return "", errUnauthorized
	}
	return user.FirstName + " " + user.LastName, nil
}
