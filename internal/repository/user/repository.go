package user

import "github.com/razedwell/go-hand/internal/model"

type Repository interface {
	FindUserByEmail(email string) (*model.User, error)
	FindUserById(id int64) (*model.User, error)
}
