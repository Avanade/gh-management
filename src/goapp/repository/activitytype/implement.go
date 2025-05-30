package activitytype

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type activityTypeRepository struct {
	*db.Database
}

func NewActivityTypeRepository(db *db.Database) ActivityTypeRepository {
	return &activityTypeRepository{db}
}

func (r *activityTypeRepository) Insert(activityType *model.ActivityType) (*model.ActivityType, error) {
	result, err := r.QueryRow("[dbo].[usp_ActivityType_Insert]",
		sql.Named("Name", activityType.Name))
	if err != nil {
		return nil, err
	}
	err = result.Scan(
		&activityType.ID)
	if err != nil {
		return nil, err
	}
	return activityType, nil
}

func (r *activityTypeRepository) Select() ([]model.ActivityType, error) {
	var activityTypes []model.ActivityType
	rows, err := r.Query("[dbo].[usp_ActivityType_Select]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var activityType model.ActivityType

		activityType.ID = v["Id"].(int64)
		activityType.Name = v["Name"].(string)

		activityTypes = append(activityTypes, activityType)
	}

	return activityTypes, nil
}
