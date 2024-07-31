package contributionarea

import (
	"errors"
	"main/model"
	repositoryContributionArea "main/repository/contributionarea"
	"strconv"
)

type contributionAreaService struct {
	repositoryContributionArea repositoryContributionArea.ContributionAreaRepository
}

// Create implements ContributionAreaService.
func (s *contributionAreaService) Create(contributionArea *model.ContributionArea) (*model.ContributionArea, error) {
	return s.repositoryContributionArea.Create(contributionArea)
}

// Get implements ContributionAreaService.
func (s *contributionAreaService) Get(offset, filter, orderby, ordertype, search string) (contributionAreas []model.ContributionArea, total int, err error) {
	if search != "" || (offset != "" && filter != "") {
		filterInt, err := strconv.Atoi(filter)
		if err != nil {
			return nil, 0, err
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, 0, err
		}
		contributionAreas, err = s.repositoryContributionArea.GetByOption(offsetInt, filterInt, orderby, ordertype, search)
		if err != nil {
			return nil, 0, err
		}
	} else {
		contributionAreas, err = s.repositoryContributionArea.GetAll()
		if err != nil {
			return nil, 0, err
		}
	}
	total, err = s.repositoryContributionArea.GetTotal()
	if err != nil {
		return nil, 0, err
	}
	return contributionAreas, total, nil
}

// GetByID implements ContributionAreaService.
func (s *contributionAreaService) GetByID(id string) (*model.ContributionArea, error) {
	parseId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repositoryContributionArea.GetByID(parseId)
}

// Update implements ContributionAreaService.
func (s *contributionAreaService) Update(id string, contributionArea *model.ContributionArea) (*model.ContributionArea, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repositoryContributionArea.Update(parsedId, contributionArea)
}

// Validate implements ContributionAreaService.
func (s *contributionAreaService) Validate(contributionArea *model.ContributionArea) error {
	if contributionArea == nil {
		return errors.New("contribution area is empty")
	}
	if contributionArea.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func NewContributionAreaService(repositoryContributionArea repositoryContributionArea.ContributionAreaRepository) ContributionAreaService {
	return &contributionAreaService{repositoryContributionArea}
}
