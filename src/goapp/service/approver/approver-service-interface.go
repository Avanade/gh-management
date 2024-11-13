package approver

import (
	"main/model"
)

type ApproverService interface {
	GetApproversByApprovalTypeId(approvalTypeId int) ([]model.Approver, error)
}
