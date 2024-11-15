package approvaltype

import (
	"main/model"
)

type ApprovalTypeService interface {
	Get(opt *model.FilterOptions) ([]model.ApprovalType, int64, error)
	GetById(id int) (*model.ApprovalType, error)
}
