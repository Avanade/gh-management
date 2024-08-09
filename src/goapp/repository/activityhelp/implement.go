package activityhelp

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type activityHelpRepository struct {
	db.Database
}

func NewActivityHelpRepository(database db.Database) ActivityHelpRepository {
	return &activityHelpRepository{database}
}

func (r *activityHelpRepository) Insert(activityId int, helpTypeId int, details string) (*model.ActivityHelp, error) {
	row, err := r.QueryRow("[dbo].[usp_CommunityActivityHelpType_Insert]",
		sql.Named("ActivityId", activityId),
		sql.Named("HelpTypeId", helpTypeId),
		sql.Named("Details", details))
	if err != nil {
		return nil, err
	}
	var activityHelp model.ActivityHelp
	err = row.Scan(
		&activityHelp.ID)
	if err != nil {
		return nil, err
	}
	return &activityHelp, nil
}
