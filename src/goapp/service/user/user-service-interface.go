package user

import (
	"main/model"
)

type UserService interface {
	Create(user *model.User) error
}
