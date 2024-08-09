package model

type ActivityHelp struct {
	ID                  int64  `json:"id"`
	CommunityActivityId int64  `json:"communityActivityId"`
	HelpTypeId          int64  `json:"helpTypeId"`
	Details             string `json:"details"`
}
