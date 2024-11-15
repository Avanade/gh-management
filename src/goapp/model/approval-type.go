package model

import "time"

type ApprovalType struct {
	Id         int                  `json:"id"`
	Name       string               `json:"name"`
	Approvers  []RepositoryApprover `json:"approvers"`
	IsActive   bool                 `json:"isActive"`
	IsArchived bool                 `json:"isArchived"`
	Created    time.Time            `json:"created"`
	CreatedBy  string               `json:"createdBy"`
	Modified   time.Time            `json:"modified"`
	ModifiedBy string               `json:"modifiedBy"`
}
