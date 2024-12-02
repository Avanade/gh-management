package user

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type userRepository struct {
	*db.Database
}

func NewUserRepository(db *db.Database) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Insert(user *model.User) error {
	err := r.Execute("usp_User_Insert",
		sql.Named("UserPrincipalName", user.UserPrincipalName),
		sql.Named("Name", user.Name),
		sql.Named("GivenName", user.GivenName),
		sql.Named("Surname", user.Surname),
		sql.Named("JobTitle", user.JobTitle),
	)
	if err != nil {
		return err
	}

	return nil
}
