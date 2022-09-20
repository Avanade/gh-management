package models

type TypBasicRepo struct {
	Name        string
	Requestor   string
	Description string
}

type TypRequestedRepoSummary struct {
	Date         string
	Organization string
	Repos        []TypBasicRepo
}
