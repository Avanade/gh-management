package approvalType

import (
	"main/model"
)

type ApprovalTypeRepository interface {
	Insert(*model.ApprovalType) (int, error)
	Select() ([]model.ApprovalType, error)
	SelectById(id int) (*model.ApprovalType, error)
	SelectByOption(opt model.FilterOptions) ([]model.ApprovalType, error)
	Total() (int64, error)
}
