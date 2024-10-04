package activitycontributionarea

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type activityContributionAreaRepository struct {
	db.Database
}

func NewActivityContributionAreaRepository(database db.Database) ActivityContributionAreaRepository {
	return &activityContributionAreaRepository{database}
}

func (r *activityContributionAreaRepository) Insert(activityContributionArea *model.ActivityContributionArea) (*model.ActivityContributionArea, error) {
	result, err := r.QueryRow("[dbo].[usp_CommunityActivityContributionArea_Insert]",
		sql.Named("CommunityActivityId", activityContributionArea.ActivityId),
		sql.Named("ContributionAreaId", activityContributionArea.ContributionAreaId),
		sql.Named("IsPrimary", activityContributionArea.IsPrimary))
	if err != nil {
		return nil, err
	}
	err = result.Scan(
		&activityContributionArea.ID)
	if err != nil {
		return nil, err
	}
	return activityContributionArea, nil
}

func (r *activityContributionAreaRepository) SelectByActivityId(activityId int64) ([]model.ActivityContributionArea, error) {
	var activityContributionAreas []model.ActivityContributionArea
	row, err := r.Query("[dbo].[usp_CommunityActivityContributionArea_Select_ByActivityId]",
		sql.Named("ActivityId", activityId))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	mapRows, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var activityContributionArea model.ActivityContributionArea

		activityContributionArea.ID = v["Id"].(int64)
		activityContributionArea.ActivityId = v["CommunityActivityId"].(int64)
		activityContributionArea.IsPrimary = v["IsPrimary"].(bool)

		activityContributionArea.ContributionAreaId = v["ContributionAreaId"].(int64)
		activityContributionArea.ContributionArea.Name = v["ContributionAreaName"].(string)

		activityContributionAreas = append(activityContributionAreas, activityContributionArea)
	}

	return activityContributionAreas, nil
}
