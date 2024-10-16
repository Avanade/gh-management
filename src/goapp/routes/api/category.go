package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"main/pkg/session"

	"github.com/gorilla/mux"
)

type CategoryDto struct {
	Id               int                   `json:"id"`
	Name             string                `json:"name"`
	Created          string                `json:"created"`
	CreatedBy        string                `json:"createdBy"`
	Modified         string                `json:"modified"`
	ModifiedBy       string                `json:"modifiedBy"`
	CategoryArticles []CategoryArticlesDto `json:"categoryArticles"`
}

type CategoryArticlesDto struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Url          string `json:"Url"`
	Body         string `json:"Body"`
	CategoryId   int    `json:"CategoryId"`
	CategoryName string `json:"CategoryName"`
	Created      string `json:"created"`
	CreatedBy    string `json:"createdBy"`
	Modified     string `json:"modified"`
	ModifiedBy   string `json:"modifiedBy"`
}

func CategoryAPIHandler(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body CategoryDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{

		"Name":       body.Name,
		"CreatedBy":  username,
		"ModifiedBy": username,
		"Id":         body.Id,
	}

	result, err := db.CategoryInsert(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, c := range body.CategoryArticles {

		categoryArticles := map[string]interface{}{

			"Id":          0,
			"Name ":       c.Name,
			"Url":         c.Url,
			"Body":        c.Body,
			"CategoryId ": id,
			"CreatedBy":   username,
			"ModifiedBy":  username,
		}

		_, err := db.CategoryArticlesInsert(categoryArticles)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CategoryListAPIHandler(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get project list
	communities, err := db.CategorySelect()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communities)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCategoryArticlesByCategoryId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	categoryArticles, err := db.CategoryArticlesselectById(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(categoryArticles)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCategoryArticlesById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id

	categoryArticles, err := db.CategoryArticlesSelectByArticlesID(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(categoryArticles)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	category, err := db.CategorySelectById(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(category)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func UpdateCategoryArticlesById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body CategoryArticlesDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := map[string]interface{}{

		"Name":       body.CategoryName,
		"CreatedBy":  username,
		"ModifiedBy": username,
		"Id":         body.CategoryId,
	}

	result, err := db.CategoryInsert(params)
	if err != nil {
		logger.LogException(err)
	}

	id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	param := map[string]interface{}{
		"Id":         body.Id,
		"Name":       body.Name,
		"Url":        body.Url,
		"Body":       body.Body,
		"CategoryId": id,
		"CreatedBy":  username,
		"ModifiedBy": username,
	}

	err = db.CategoryArticlesUpdate(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	var body CategoryDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := map[string]interface{}{
		"Name":       body.Name,
		"CreatedBy":  username,
		"ModifiedBy": username,
		"Id":         body.Id,
	}

	_, err = db.CategoryInsert(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
