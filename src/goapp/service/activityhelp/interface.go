package activityhelp

import "main/model"

type ActivityHelpService interface {
	Validate(activityId, helpTypeId int, details string) error
	Insert(activityId, helpTypeId int, details string) (*model.ActivityHelp, error)
}
