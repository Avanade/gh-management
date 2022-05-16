package models

type TypPageData struct {
	Header    interface{}
	Profile   interface{}
	ProfileGH TypGitHubUser
	Content   interface{}
}

type TypGitHubUser struct {
	LoggedIn    bool
	Username    string `json:"login"`
	NodeId      string `json:"node_id"`
	AvatarUrl   string `json:"avatar_url"`
	AccessToken string
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
