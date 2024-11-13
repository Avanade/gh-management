package model

type ApprovalType struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Approvers  []Approver `json:"approvers"`
	IsActive   bool       `json:"isActive"`
	IsArchived bool       `json:"isArchived"`
}
