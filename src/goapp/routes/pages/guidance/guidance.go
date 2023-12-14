package routes

import (
	"log"
	"net/http"

	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

func GuidanceHandler(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	id := req["id"]

	// check ung rout is has admin  ias
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Id":      id,
		"IsAdmin": isAdmin,
	}
	template.UseTemplate(&w, r, "/guidance/index", data)
}

func NewArticleHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "/guidance/article/new", nil)
}

func EditCategoryHandler(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	id := req["id"]
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		log.Println(err.Error())
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
	template.UseTemplate(&w, r, "/guidance/category/edit", data)
}

func EditArticleHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		log.Println(err.Error())
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
	template.UseTemplate(&w, r, "/guidance/article/edit", data)
}
