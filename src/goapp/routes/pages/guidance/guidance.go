package routes

import (

	//session "main/pkg/session"
	template "main/pkg/template"
	"net/http"

	"github.com/gorilla/mux"
	//models "main/models"
)

func GuidanceHandler(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	id := req["id"]

	data := map[string]interface{}{
		"Id": id,
	}

	template.UseTemplate(&w, r, "/guidance/guidance", data)
}
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	//guidances := [{"test"}]
	//guidances.append("test2")
	template.UseTemplate(&w, r, "/guidance/Categories", nil)
}
