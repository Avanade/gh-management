package models

type TypPageData struct {
	Header  interface{}
	Profile interface{}
	Content interface{}
}

type TypGitHubUser struct {
	LoggedIn    bool
	Username    string `json:"login"`
	NodeId      string `json:"node_id"`
	AvatarUrl   string `json:"avatar_url"`
	AccessToken string
}

type TypHeaders struct {
	Menu []TypMenu
}

type TypMenu struct {
	Name string
	Url  string
}
