package approvaltype

import (
	"main/model"
)

type ApprovalTypeService interface {
	GetApprovalTypes(opt *model.FilterOptions) ([]model.ApprovalType, error)
	GetTotalApprovalTypes() (int64, error)
}
