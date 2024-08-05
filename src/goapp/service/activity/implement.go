package activity

import (
	"errors"
	"main/model"
	repositoryActivity "main/repository/activity"
	"strconv"
)

type activityService struct {
	repositoryActivity repositoryActivity.ActivityRepository
}

// Create implements ActivityService.
func (s *activityService) Create(activity *model.Activity) (*model.Activity, error) {
	return s.repositoryActivity.Insert(activity)
}

// GetAll implements ActivityService.
func (s *activityService) Get(offset, filter, orderby, ordertype, search string) (activities []model.Activity, total int64, err error) {
	if search != "" || (offset != "" && filter != "") {
		filterInt, err := strconv.ParseInt(filter, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		offsetInt, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return nil, 0, err
		}
		createdBy := "USER"
		total, err = s.repositoryActivity.TotalByOptions(search, createdBy)
		if err != nil {
			return nil, 0, err
		}
		activities, err = s.repositoryActivity.SelectByOptions(offsetInt, filterInt, orderby, ordertype, search, createdBy)
		if err != nil {
			return nil, 0, err
		}
	} else {
		total, err = s.repositoryActivity.Total()
		if err != nil {
			return nil, 0, err
		}
		activities, err = s.repositoryActivity.Select()
		if err != nil {
			return nil, 0, err
		}
	}

	if err != nil {
		return nil, 0, err
	}
	return activities, total, nil
}

// GetById implements ActivityService.
func (s *activityService) GetById(id string) (*model.Activity, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repositoryActivity.SelectById(parsedId)
}

// Validate implements ActivityService.
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
	if activity.ActivityTypeId == 0 {
		return errors.New("activity type id is required")
	}
	return nil
}

func NewActivityService(activityRepository repositoryActivity.ActivityRepository) ActivityService {
	return &activityService{
		repositoryActivity: activityRepository,
	}
}
