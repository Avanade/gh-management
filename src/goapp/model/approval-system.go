package model

type ApprovalSystemRequestBody struct {
	ApplicationId       string
	ApplicationModuleId string
	Emails              []string
	Subject             string
	Body                string
	RequesterEmail      string
}

type ApprovalSystemResponse struct {
	ItemId string `json:"itemId"`
}
