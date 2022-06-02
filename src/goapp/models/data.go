package models

type TypPageData struct {
	Header    interface{}
	Profile   interface{}
	ProfileGH TypGitHubUser
	Content   interface{}
}

type TypGitHubUser struct {
	LoggedIn    bool
	Id          int    `json:"id"`
	Username    string `json:"login"`
	NodeId      string `json:"node_id"`
	AvatarUrl   string `json:"avatar_url"`
	AccessToken string
	IsValid     bool
}

type TypHeaders struct {
	Menu          []TypMenu
	ExternalLinks []TypMenu
	Page          string
}

type TypMenu struct {
	Name     string
	Url      string
	IconPath string
}

type TypNewProjectReqBody struct {
	Name             string `json:"name"`
	Coowner          string `json:"coowner"`
	Description      string `json:"description"`
	ConfirmAvaIP     bool   `json:"confirmAvaIP"`
	ConfirmSecIPScan bool   `json:"confirmSecIPScan"`
}

type TypUpdateApprovalStatusReqBody struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"Remarks"`
	ResponseDate string `json:"responseDate"`
}
