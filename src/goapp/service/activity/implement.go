package activity

import (
	"errors"
	"main/model"
	"main/repository"
	"strconv"
)

type activityService struct {
	*repository.Repository
}

func NewActivityService(repo *repository.Repository) ActivityService {
	return &activityService{
		Repository: repo,
	}
}

func (s *activityService) Create(activity *model.Activity) (*model.Activity, error) {
	if activity.ActivityType.ID == 0 {
		activityType, err := s.Repository.ActivityType.Insert(&activity.ActivityType)
		if err != nil {
			return nil, err
		}
		activity.ActivityTypeId = activityType.ID
		activity.ActivityType = model.ActivityType{
			ID:   activityType.ID,
			Name: activityType.Name,
		}
	}

	activity, err := s.Repository.Activity.Insert(activity)
	if err != nil {
		return nil, err
	}

	for _, contributionArea := range activity.ActivityContributionAreas {
		contributionArea.ActivityId = activity.ID
		_, err := s.Repository.ActivityContributionArea.Insert(&contributionArea)
		if err != nil {
			return nil, err
		}
	}

	return activity, nil
}

func (s *activityService) Get(offset, filter, orderby, ordertype, search, createdBy string) (activities []model.Activity, total int64, err error) {
	if search != "" || (offset != "" && filter != "") {
		filterInt, err := strconv.ParseInt(filter, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		offsetInt, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		total, err = s.Repository.Activity.TotalByOptions(search, createdBy)
		if err != nil {
			return nil, 0, err
		}
		activities, err = s.Repository.Activity.SelectByOptions(offsetInt, filterInt, orderby, ordertype, search, createdBy)
		if err != nil {
			return nil, 0, err
		}
	} else {
		total, err = s.Repository.Activity.Total()
		if err != nil {
			return nil, 0, err
		}
		activities, err = s.Repository.Activity.Select()
		if err != nil {
			return nil, 0, err
		}
	}

	if err != nil {
		return nil, 0, err
	}
	return activities, total, nil
}

func (s *activityService) GetById(id string) (*model.Activity, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.Repository.Activity.SelectById(parsedId)
}

func (s *activityService) Validate(activity *model.Activity) error {
	if activity == nil {
		return errors.New("activity is empty")
	}
	if activity.CommunityId == 0 {
		return errors.New("community id is required")
	}
	if activity.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
