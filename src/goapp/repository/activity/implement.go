package activity

import (
	"database/sql"
	"main/model"
	"main/repository"
	"time"
)

type activityRepository struct {
	repository.Database
}

// Insert implements ActivityRepository.
func (r *activityRepository) Insert(activity *model.Activity) (*model.Activity, error) {
	result, err := r.QueryRow("[dbo].[usp_CommunityActivity_Insert]",
		sql.Named("CommunityId", activity.CommunityId),
		sql.Named("Name", activity.Name),
		sql.Named("ActivityTypeId", activity.ActivityTypeId),
		sql.Named("Url", activity.Url),
		sql.Named("Date", activity.Date),
		sql.Named("CreatedBy", activity.CreatedBy))
	if err != nil {
		return nil, err
	}
	err = result.Scan(
		&activity.ID,
		&activity.Created)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// Select implements ActivityRepository.
func (r *activityRepository) Select() ([]model.Activity, error) {
	var activities []model.Activity
	row, err := r.Query("[dbo].[usp_CommunityActivity_Select]")
	if err != nil {
		return nil, err
	}
	defer row.Close()

	mapRows, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var activity model.Activity

		activity.ID = v["Id"].(int64)
		activity.CommunityId = v["CommunityId"].(int64)
		activity.Name = v["Name"].(string)
		activity.ActivityTypeId = v["ActivityTypeId"].(int64)
		activity.Url = v["Url"].(string)
		activity.Date = v["Date"].(time.Time)
		activity.Created = v["Created"].(time.Time)
		activity.CreatedBy = v["CreatedBy"].(string)
		activity.Modified = v["Modified"].(time.Time)
		if v["ModifiedBy"] != nil {
			activity.ModifiedBy = v["ModifiedBy"].(string)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// SelectById implements ActivityRepository.
func (r *activityRepository) SelectById(id int64) (*model.Activity, error) {
	var activity model.Activity
	row, err := r.QueryRow("[dbo].[usp_CommunityActivity_Select_ById]",
		sql.Named("Id", id))
	if err != nil {
		return nil, err
	}

	var modifiedBy sql.NullString
	var url sql.NullString

	err = row.Scan(
		&activity.ID,
		&activity.Name,
		&activity.CommunityId,
		&activity.CommunityName,
		&activity.ActivityTypeId,
		&activity.ActivityTypeName,
		&activity.PrimaryContributionAreaId,
		&activity.PrimaryContributionAreaName,
		&url,
		&activity.Date,
		&activity.Created,
		&activity.CreatedBy,
		&activity.Modified,
		&modifiedBy,
	)
	if err != nil {
		return nil, err
	}

	if url.Valid {
		activity.Url = url.String
	}
	if modifiedBy.Valid {
		activity.ModifiedBy = modifiedBy.String
	}

	return &activity, nil
}

// SelectByOptions implements ActivityRepository.
func (r *activityRepository) SelectByOptions(offset int64, filter int64, orderBy string, orderType string, search string, createdBy string) ([]model.Activity, error) {
	var activities []model.Activity
	rows, err := r.Query("[dbo].[usp_CommunityActivity_Select_ByOptionAndCreatedBy]",
		sql.Named("Offset", offset),
		sql.Named("Filter", filter),
		sql.Named("Search", search),
		sql.Named("OrderBy", orderBy),
		sql.Named("OrderType", orderType),
		sql.Named("CreatedBy", createdBy))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var activity model.Activity

		activity.ID = v["Id"].(int64)
		activity.CommunityId = v["CommunityId"].(int64)
		activity.Name = v["Name"].(string)
		activity.ActivityTypeId = v["ActivityTypeId"].(int64)
		activity.Url = v["Url"].(string)
		activity.Date = v["Date"].(time.Time)
		activity.Created = v["Created"].(time.Time)
		activity.CreatedBy = v["CreatedBy"].(string)
		activity.Modified = v["Modified"].(time.Time)
		if v["ModifiedBy"] != nil {
			activity.ModifiedBy = v["ModifiedBy"].(string)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// Total implements ActivityRepository.
func (r *activityRepository) Total() (int64, error) {
	var total int64
	row, err := r.QueryRow("[dbo].[usp_CommunityActivity_TotalCount]")
	if err != nil {
		return 0, err
	}
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// TotalByOptions implements ActivityRepository.
func (r *activityRepository) TotalByOptions(search string, createdBy string) (int64, error) {
	var total int64
	row, err := r.QueryRow("[dbo].[usp_CommunityActivity_TotalCount_ByOptionAndCreatedBy]",
		sql.Named("Search", search),
		sql.Named("CreatedBy", createdBy))
	if err != nil {
		return 0, err
	}
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func NewActivityRepository(db repository.Database) ActivityRepository {
	return &activityRepository{db}
}
