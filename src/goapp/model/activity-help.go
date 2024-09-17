package model

type ActivityHelp struct {
	ID         int64  `json:"id"`
	ActivityId int64  `json:"activityId"`
	HelpTypeId int64  `json:"helpTypeId"`
	Details    string `json:"details"`
}
