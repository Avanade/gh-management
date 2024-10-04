package contributionarea

import (
	"errors"
	"main/model"
	"main/repository"
	"strconv"
)

type contributionAreaService struct {
	Repository *repository.Repository
}

func NewContributionAreaService(repo *repository.Repository) ContributionAreaService {
	return &contributionAreaService{
		Repository: repo,
	}
}

func (s *contributionAreaService) Create(contributionArea *model.ContributionArea) (*model.ContributionArea, error) {
	return s.Repository.ContributionArea.Insert(contributionArea)
}

func (s *contributionAreaService) Get(offset, filter, orderby, ordertype, search string) (contributionAreas []model.ContributionArea, total int64, err error) {
	if search != "" || (offset != "" && filter != "") {
		filterInt, err := strconv.Atoi(filter)
		if err != nil {
			return nil, 0, err
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, 0, err
		}
		contributionAreas, err = s.Repository.ContributionArea.SelectByOption(offsetInt, filterInt, orderby, ordertype, search)
		if err != nil {
			return nil, 0, err
		}
	} else {
		contributionAreas, err = s.Repository.ContributionArea.Select()
		if err != nil {
			return nil, 0, err
		}
	}
	total, err = s.Repository.ContributionArea.Total()
	if err != nil {
		return nil, 0, err
	}
	return contributionAreas, total, nil
}

func (s *contributionAreaService) GetByID(id string) (*model.ContributionArea, error) {
	parseId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.Repository.ContributionArea.SelectById(parseId)
}

func (s *contributionAreaService) Update(id string, contributionArea *model.ContributionArea) (*model.ContributionArea, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.Repository.ContributionArea.Update(parsedId, contributionArea)
}

func (s *contributionAreaService) Validate(contributionArea *model.ContributionArea) error {
	if contributionArea == nil {
		return errors.New("contribution area is empty")
	}
	if contributionArea.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
