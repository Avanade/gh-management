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

type TypCommunity struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Url          string   `json:"url"`
	Description  string   `json:"description"`
	Notes        string   `json:"notes"`
	TradeAssocId string   `json:"tradeAssocId"`
	Created      string   `json:"created"`
	CreatedBy    string   `json:"createdBy"`
	Modified     string   `json:"modified"`
	ModifiedBy   string   `json:"modifiedBy"`
	Sponsors     []string `json:"sponsors"`
}
type TypCommunitySponsors struct {
	Id                string `json:"id"`
	CommunityId       string `json:"communityId"`
	UserPrincipalName string `json:"userprincipalname"`
	Created           string `json:"created"`
	CreatedBy         string `json:"createdBy"`
	Modified          string `json:"modified"`
	ModifiedBy        string `json:"modifiedBy"`
}
