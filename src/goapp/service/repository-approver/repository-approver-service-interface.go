package repositoryApprover

import (
	"main/model"
)

type RepositoryApproverService interface {
	Create(approver *model.RepositoryApprover) error
	Get(approvalTypeId int) ([]model.RepositoryApprover, error)
}
