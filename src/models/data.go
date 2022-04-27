package models

type TypPageData struct {
	Header  interface{}
	Profile interface{}
	Content interface{}
}

type TypHeaders struct {
	Menu []TypMenu
}

type TypMenu struct {
	Name string
	Url  string
}
