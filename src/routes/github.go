package routes

import (
	"main/models"
	template "main/pkg/template"
	"net/http"
)

func GithubHandler(w http.ResponseWriter, r *http.Request, data *models.TypPageData) {
	var users []typUsers
	users = append(users, typUsers{id: 1, name: "Rainier"})
	users = append(users, typUsers{id: 2, name: "Isko"})
	users = append(users, typUsers{id: 3, name: "Dennis"})
	users = append(users, typUsers{id: 4, name: "Jerrico"})
	users = append(users, typUsers{id: 5, name: "Daryl"})

	var items []typItems
	items = append(items, typItems{category: "Pencil", name: "Mongol #1"})
	items = append(items, typItems{category: "Pencil", name: "Mongol #2"})
	items = append(items, typItems{category: "Ballpen", name: "Panda Black"})

	data.Content = typContent{users: users, items: items}
	tmpl := template.UseTemplate("github")
	tmpl.Execute(w, data)
}

type typContent struct {
	users []typUsers
	items []typItems
}
type typUsers struct {
	id   int
	name string
}

type typItems struct {
	category string
	name     string
}
