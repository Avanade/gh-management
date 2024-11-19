package user

import (
	"main/model"
	"main/repository"
)

type userService struct {
	Repository *repository.Repository
}

func NewUserService(repository *repository.Repository) UserService {
	return &userService{repository}
}

func (s *userService) Create(user *model.User) error {
	return s.Repository.User.Insert(user)
}
