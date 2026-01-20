package user

import (
	"context"

	"github.com/razedwell/go-hand/internal/model"
)

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserById(ctx context.Context, id int64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}
