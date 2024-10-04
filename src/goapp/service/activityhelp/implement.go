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

func (s *activityHelpService) Create(activityHelp *model.ActivityHelp) error {
	return s.Repository.ActivityHelp.Insert(activityHelp.ActivityId, activityHelp.HelpTypeId, activityHelp.Details)
}

func (s *activityHelpService) Validate(activityHelp *model.ActivityHelp) error {
	if activityHelp.ActivityId == 0 {
		return errors.New("no activity id provided")
	}
	if activityHelp.HelpTypeId == 0 {
		return errors.New("help type is required")
	}
	if activityHelp.Details == "" {
		return errors.New("details is required")
	}
	return nil
}
