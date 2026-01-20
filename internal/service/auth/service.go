package auth

import (
	"context"
	"errors"

	"github.com/razedwell/go-hand/internal/repository/user"
	"github.com/razedwell/go-hand/internal/security"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var errUnauthorized = errors.New("unauthorized")

type Service struct {
	users user.Repository
	jwt   *security.JWTManager
}

func NewService(users user.Repository, jwt *security.JWTManager) *Service {
	return &Service{users, jwt}
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.users.FindUserByEmail(ctx, email)
	if err != nil || security.VerifyPassword(user.PasswordHash, password) == false {
		return "", errUnauthorized
	}
	return s.jwt.Generate(user.ID)
}
