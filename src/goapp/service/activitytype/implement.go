package activitytype

import (
	"main/model"
	repositoryActivityType "main/repository/activitytype"
)

type activityTypeService struct {
	activityTypeRepository repositoryActivityType.ActivityTypeRepository
}

func NewActivityTypeService(activityTypeRepository repositoryActivityType.ActivityTypeRepository) ActivityTypeService {
	return &activityTypeService{activityTypeRepository}
}

func (s *activityTypeService) Get() ([]model.ActivityType, error) {
	return s.activityTypeRepository.Select()
}
