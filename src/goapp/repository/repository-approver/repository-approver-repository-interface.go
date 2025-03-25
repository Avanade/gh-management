package repositoryApprover

import (
	"main/model"
)

type RepositoryApproverRepository interface {
	Insert(approver *model.RepositoryApprover) error
	SelectByApprovalTypeId(approvalTypeId int) ([]model.RepositoryApprover, error)
}
