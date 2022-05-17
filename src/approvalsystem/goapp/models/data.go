package models

type TypPageData struct {
	Header  interface{}
	Profile interface{}
	Content interface{}
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
