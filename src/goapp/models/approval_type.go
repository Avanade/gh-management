package models

type ApprovalType struct {
	Id                        int
	Name                      string
	ApproverUserPrincipalName string
	IsActive                  bool
}
