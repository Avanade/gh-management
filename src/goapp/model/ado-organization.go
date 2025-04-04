package model

import "time"

type AdoOrganizationRequest struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Purpose   string    `json:"purpose"`
	Created   time.Time `json:"created"`
	CreatedBy string    `json:"createdBy"`
}

type ApprovalRequest struct {
	Id                    int64     `json:"id"`
	ApprovalDate          time.Time `json:"approvalDate"`
	ApprovalDescription   string    `json:"approvalDescription"`
	ApprovalRemarks       string    `json:"approvalRemarks"`
	ApprovalStatus        string    `json:"approvalStatus"`
	ApproverPrincipalName string    `json:"approverPrincipalName"`
}
