package approver

import (
	"main/model"
)

type ApproverService interface {
	Get(approvalTypeId int) ([]model.RepositoryApprover, error)
}
