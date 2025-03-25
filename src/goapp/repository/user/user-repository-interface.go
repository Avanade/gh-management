package user

import (
	"main/model"
)

type UserRepository interface {
	Insert(*model.User) error
}
