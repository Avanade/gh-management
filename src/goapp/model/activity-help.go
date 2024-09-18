package model

type ActivityHelp struct {
	ActivityId int64  `json:"activityId"`
	HelpTypeId int64  `json:"helpTypeId"`
	Details    string `json:"details"`

	HelpType HelpType `json:"helpType"`
}
