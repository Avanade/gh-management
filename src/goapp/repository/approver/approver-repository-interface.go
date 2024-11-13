package approver

import (
	"main/model"
)

type ApproverRepository interface {
	GetApproversByApprovalTypeId(approvalTypeId int) ([]model.Approver, error)
}
