package routes

import (
	"encoding/json"
	"log"
	"main/pkg/msgraph"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllUserFromActiveDirectory(w http.ResponseWriter, r *http.Request) {
	users, err := msgraph.GetAllUsers()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func SearchUserFromActiveDirectory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]

	users, err := msgraph.SearchUsers(search)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
