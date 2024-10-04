package activityhelp

type ActivityHelpRepository interface {
	Insert(activityId, helpTypeId int64, details string) error
}
