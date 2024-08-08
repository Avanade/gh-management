package activitytype

import (
	"main/model"
	"main/repository"
)

type activityTypeService struct {
	Repository *repository.Repository
}

func NewActivityTypeService(repo *repository.Repository) ActivityTypeService {
	return &activityTypeService{Repository: repo}
}

func (s *activityTypeService) Get() ([]model.ActivityType, error) {
	return s.Repository.ActivityType.Select()
}
