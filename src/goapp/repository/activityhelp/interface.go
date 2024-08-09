package activityhelp

import "main/model"

type ActivityHelpRepository interface {
	Insert(activityId, helpTypeId int, details string) (*model.ActivityHelp, error)
}
