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

type TypProjectApprovals struct {
	Id                         int64
	ProjectId                  int64
	ProjectName                string
	ProjectCoowner             string
	ProjectDescription         string
	RequesterName              string
	RequesterGivenName         string
	RequesterSurName           string
	RequesterUserPrincipalName string
	CoownerName                string
	CoownerGivenName           string
	CoownerSurName             string
	CoownerUserPrincipalName   string
	ApprovalTypeId             int64
	ApprovalType               string
	ApproverUserPrincipalName  string
	ApprovalDescription        string
}

type TypApprovalSystemPost struct {
	ApplicationId       string
	ApplicationModuleId string
	Email               string
	Subject             string
	Body                string
	RequesterEmail      string
}

type TypApprovalSystemPostResponse struct {
	ItemId string `json:"itemId"`
}

type TypUpdateApprovalStatusReqBody struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"Remarks"`
	ResponseDate string `json:"responseDate"`
}
