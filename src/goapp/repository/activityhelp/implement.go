package activityhelp

import (
	"database/sql"
	db "main/infrastructure/database"
)

type activityHelpRepository struct {
	db.Database
}

func NewActivityHelpRepository(database db.Database) ActivityHelpRepository {
	return &activityHelpRepository{database}
}

func (r *activityHelpRepository) Insert(activityId int64, helpTypeId int64, details string) error {
	err := r.Execute("[dbo].[usp_CommunityActivityHelpType_Insert]",
		sql.Named("ActivityId", activityId),
		sql.Named("HelpTypeId", helpTypeId),
		sql.Named("Details", details))
	if err != nil {
		return err
	}
	return nil
}
