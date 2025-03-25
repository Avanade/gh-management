package model

type User struct {
	UserPrincipalName string `json:"userPrincipalName"`
	Name              string `json:"name"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	JobTitle          string `json:"jobTitle"`
}
