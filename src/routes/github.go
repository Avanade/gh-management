package routes

import (
	template "main/pkg/template"
	"net/http"
)

func GithubHandler(w http.ResponseWriter, r *http.Request) {
	var rules []TypGHMessage
	rules = append(rules, TypGHMessage{Title: "Rule #1", Message: "The 4th parameter on UseTemplate can be nil."})
	rules = append(rules, TypGHMessage{Title: "Rule #2", Message: "Parameter can be of any type."})
	rules = append(rules, TypGHMessage{Title: "Rule #3", Message: "If the data uses a struct, it must be exported."})
	rules = append(rules, TypGHMessage{Title: "Rule #4", Message: "If you need to pass-on multiple data, create a struct that includes all other struct."})
	template.UseTemplate(&w, r, "github", rules)
}

type TypGHMessage struct {
	Title   string
	Message string
}
