package model

type RepositoryApprover struct {
	ApprovalTypeId int    `json:"approvalTypeId"`
	ApproverEmail  string `json:"approverEmail"`
	ApproverName   string `json:"approverName"`
}
