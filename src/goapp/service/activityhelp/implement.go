package activityhelp

import (
	"errors"
	"main/model"
	"main/repository"
)

type activityHelpService struct {
	*repository.Repository
}

func NewActivityHelpService(repo *repository.Repository) ActivityHelpService {
	return &activityHelpService{
		Repository: repo,
	}
}

func (s *activityHelpService) Insert(activityId, helpTypeId int, details string) (*model.ActivityHelp, error) {
	return s.Repository.ActivityHelp.Insert(activityId, helpTypeId, details)
}

func (s *activityHelpService) Validate(activityId, helpTypeId int, details string) error {
	if activityId == 0 {
		return errors.New("no activity id provided")
	}
	if helpTypeId == 0 {
		return errors.New("help type is required")
	}
	if details == "" {
		return errors.New("details is required")
	}
	return nil
}
