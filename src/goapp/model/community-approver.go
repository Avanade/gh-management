package model

type CommunityApprover struct {
	Id                        int64  `json:"id"`
	ApproverUserPrincipalName string `json:"approverUserPrincipalName"`
	GuidanceCategory          string `json:"guidanceCategory"`
}
