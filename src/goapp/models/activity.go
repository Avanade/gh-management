package models

type Activity struct {
	Name        string
	Url         string
	Date        string
	TypeId      int
	CommunityId int
	CreatedBy   string
	ModifiedBy  string
}
