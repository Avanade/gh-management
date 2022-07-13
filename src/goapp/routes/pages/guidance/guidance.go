package routes

import (

	//session "main/pkg/session"
	template "main/pkg/template"
	"net/http"

	"github.com/gorilla/mux"
	//models "main/models"
	session "main/pkg/session"
)

func GuidanceHandler(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	id := req["id"]
	var IsAdmin = false
	// check ung rout is has admin  ias
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isAdmin {
		IsAdmin = true
	}

	data := map[string]interface{}{
		"Id":      id,
		"IsAdmin": IsAdmin,
	}

	template.UseTemplate(&w, r, "/guidance/guidance", data)
}
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "/guidance/Categories", nil)
}
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "Not enough privilege to do the action.", http.StatusForbidden)
		return
	}
	data := map[string]interface{}{
		"Id": id,
	}
	template.UseTemplate(&w, r, "/guidance/article", data)
}
