package models

type User struct {
	HumanName       string
	SourceControlUn string
	UserId          string
	CorpUserEmail   string
}

type UserService interface {
	User(token string) (User, error)
}
