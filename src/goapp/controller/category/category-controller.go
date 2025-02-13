package category

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/model"
	"main/pkg/appinsights_wrapper"
	"main/pkg/session"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type categoryController struct {
	*service.Service
}

func NewCategoryController(service *service.Service) CategoryController {
	return &categoryController{service}

}

func (c *categoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	categories, err := c.Service.Category.GetAll()
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (c *categoryController) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok || idStr == "" {
		logger.TrackException(errors.New("no parameters found"))
		http.Error(w, "no parameters found", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.TrackException(errors.New("invalid id parameter"))
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	category, err := c.Service.Category.GetById(id)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (c *categoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	w.Header().Set("Content-Type", "application/json")
	var categories model.Category

	err := json.NewDecoder(r.Body).Decode(&categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	categories.CreatedBy = username
	categories.ModifiedBy = username
	result, err := c.Service.Category.Insert(&categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving a new Category"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (c *categoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	if len(params) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}
	if params["id"] == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var data model.Category
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	data.CreatedBy = username
	data.ModifiedBy = username

	err = c.Service.Category.Update(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error updating the article"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
