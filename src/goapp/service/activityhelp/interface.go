package activityhelp

import "main/model"

type ActivityHelpService interface {
	Validate(activityHelp *model.ActivityHelp) error
	Create(activityHelp *model.ActivityHelp) (*model.ActivityHelp, error)
}
