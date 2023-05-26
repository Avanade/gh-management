package models

type ApprovalType struct {
	Id                        int
	Name                      string
	ApproverUserPrincipalName string
	IsActive                  bool
	IsArchived                bool
	Created                   string
	CreatedBy                 string
	Modified                  string
	ModifiedBy                string
}
