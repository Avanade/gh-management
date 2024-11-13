package approvaltype

import (
	"main/model"
	"main/repository"
)

type approvalTypeService struct {
	Repository *repository.Repository
}

func NewApprovalTypeService(repository *repository.Repository) ApprovalTypeService {
	return &approvalTypeService{repository}
}

func (s *approvalTypeService) GetApprovalTypes(opt *model.FilterOptions) ([]model.ApprovalType, error) {
	if opt == nil {
		return s.Repository.ApprovalType.GetAllApprovalTypes()
	} else {
		return s.Repository.ApprovalType.GetApprovalTypesByFilter(*opt)
	}
}

func (s *approvalTypeService) GetTotalApprovalTypes() (int64, error) {
	return s.Repository.ApprovalType.GetTotalApprovalTypes()
}
