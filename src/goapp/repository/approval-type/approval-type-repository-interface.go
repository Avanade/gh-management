package approvalType

import (
	"main/model"
)

type ApprovalTypeRepository interface {
	GetAllApprovalTypes() ([]model.ApprovalType, error)
	GetApprovalTypesByFilter(opt model.FilterOptions) ([]model.ApprovalType, error)
	GetTotalApprovalTypes() (int64, error)
}
