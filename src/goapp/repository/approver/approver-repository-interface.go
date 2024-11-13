package approver

import (
	"main/model"
)

type ApproverRepository interface {
	SelectByApprovalTypeId(approvalTypeId int) ([]model.RepositoryApprover, error)
}
