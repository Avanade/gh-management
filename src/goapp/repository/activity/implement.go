package activity

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type activityRepository struct {
	*db.Database
}

func NewActivityRepository(database *db.Database) ActivityRepository {
	return &activityRepository{database}
}

func (r *activityRepository) Insert(activity *model.Activity) (*model.Activity, error) {
	row, err := r.QueryRow("[dbo].[usp_CommunityActivity_Insert]",
		sql.Named("CommunityId", activity.CommunityId),
		sql.Named("Name", activity.Name),
		sql.Named("ActivityTypeId", activity.ActivityTypeId),
		sql.Named("Url", activity.Url),
		sql.Named("Date", activity.Date.Format("2006-01-02")),
		sql.Named("CreatedBy", activity.CreatedBy))
	if err != nil {
		return nil, err
	}
	err = row.Scan(
		&activity.ID,
		&activity.Created)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

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
		activity.Name = v["Name"].(string)
		activity.Url = v["Url"].(string)
		activity.Date = v["Date"].(time.Time)
		activity.Created = v["Created"].(time.Time)
		activity.CreatedBy = v["CreatedBy"].(string)
		if v["Modified"] != nil {
			activity.Modified = v["Modified"].(time.Time)
		}
		if v["ModifiedBy"] != nil {
			activity.ModifiedBy = v["ModifiedBy"].(string)
		}

		activity.CommunityId = v["CommunityId"].(int64)
		activity.Community = model.Community{
			ID:   v["CommunityId"].(int64),
			Name: v["CommunityName"].(string),
		}

		activity.ActivityTypeId = v["ActivityTypeId"].(int64)
		activity.ActivityType = model.ActivityType{
			ID:   v["ActivityTypeId"].(int64),
			Name: v["ActivityTypeName"].(string),
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

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
		activity.Name = v["Name"].(string)
		activity.Url = v["Url"].(string)
		activity.Date = v["Date"].(time.Time)
		activity.Created = v["Created"].(time.Time)
		activity.CreatedBy = v["CreatedBy"].(string)
		if v["Modified"] != nil {
			activity.Modified = v["Modified"].(time.Time)
		}
		if v["ModifiedBy"] != nil {
			activity.ModifiedBy = v["ModifiedBy"].(string)
		}

		activity.CommunityId = v["CommunityId"].(int64)
		activity.Community = model.Community{
			ID:   v["CommunityId"].(int64),
			Name: v["CommunityName"].(string),
		}

		activity.ActivityTypeId = v["ActivityTypeId"].(int64)
		activity.ActivityType = model.ActivityType{
			ID:   v["ActivityTypeId"].(int64),
			Name: v["ActivityTypeName"].(string),
		}

		activity.ActivityContributionAreas = append(activity.ActivityContributionAreas, model.ActivityContributionArea{
			ID:                 v["ActivityContributionAreaId"].(int64),
			ActivityId:         v["Id"].(int64),
			ContributionAreaId: v["PrimaryContributionAreaId"].(int64),
			IsPrimary:          true,
			ContributionArea: model.ContributionArea{
				ID:        v["PrimaryContributionAreaId"].(int64),
				Name:      v["PrimaryContributionAreaName"].(string),
				Created:   v["PrimaryContributionCreated"].(time.Time),
				CreatedBy: v["PrimaryContributionCreatedBy"].(string),
			},
		})

		if v["PrimaryContributionModified"] != nil {
			activity.ActivityContributionAreas[0].ContributionArea.Modified = v["PrimaryContributionModified"].(time.Time)
		}

		if v["PrimaryContributionModifiedBy"] != nil {
			activity.ActivityContributionAreas[0].ContributionArea.ModifiedBy = v["PrimaryContributionModifiedBy"].(string)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *activityRepository) SelectById(id int64) (*model.Activity, error) {
	var activity model.Activity
	rows, err := r.Query("[dbo].[usp_CommunityActivity_Select_ById]",
		sql.Named("Id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		activity.ID = v["Id"].(int64)
		activity.Name = v["Name"].(string)
		activity.Url = v["Url"].(string)
		activity.Date = v["Date"].(time.Time)
		activity.Created = v["Created"].(time.Time)
		activity.CreatedBy = v["CreatedBy"].(string)
		if v["Modified"] != nil {
			activity.Modified = v["Modified"].(time.Time)
		}
		if v["ModifiedBy"] != nil {
			activity.ModifiedBy = v["ModifiedBy"].(string)
		}

		activity.CommunityId = v["CommunityId"].(int64)
		activity.Community = model.Community{
			ID:   v["CommunityId"].(int64),
			Name: v["CommunityName"].(string),
		}

		activity.ActivityTypeId = v["ActivityTypeId"].(int64)
		activity.ActivityType = model.ActivityType{
			ID:   v["ActivityTypeId"].(int64),
			Name: v["ActivityTypeName"].(string),
		}
	}

	return &activity, nil
}

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
