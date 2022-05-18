package models

type TypRequestApproval struct {
	ApplicationId       string `json:"applicationId"`
	ApplicationModuleId string `json:"applicationModuleId"`
	Email               string `json:"email"`
	Subject             string `json:"subject"`
	Body                string `json:"body"`
}
