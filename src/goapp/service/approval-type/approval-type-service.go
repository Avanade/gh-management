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

func (s *approvalTypeService) Get(opt *model.FilterOptions) ([]model.ApprovalType, int64, error) {
	var approvalTypes []model.ApprovalType
	if opt == nil {
		data, err := s.Repository.ApprovalType.Select()
		if err != nil {
			return nil, 0, err
		}
		approvalTypes = data
	} else {
		data, err := s.Repository.ApprovalType.SelectByOption(*opt)
		if err != nil {
			return nil, 0, err
		}
		approvalTypes = data
	}

	total, err := s.Repository.ApprovalType.Total()
	if err != nil {
		return nil, 0, err
	}

	return approvalTypes, total, nil
}

func (s *approvalTypeService) GetById(id int) (*model.ApprovalType, error) {
	data, err := s.Repository.ApprovalType.SelectById(id)
	if err != nil {
		return nil, err
	}

	data.Approvers, err = s.Repository.Approver.SelectByApprovalTypeId(data.Id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
