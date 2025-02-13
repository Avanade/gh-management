package article

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

type articleController struct {
	*service.Service
}

func NewArticleController(service *service.Service) ArticleController {
	return &articleController{service}
}

func (c *articleController) GetArticlesByCategoryId(w http.ResponseWriter, r *http.Request) {
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

	category, err := c.Service.Article.GetByCategoryId(id)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (c *articleController) GetArticleById(w http.ResponseWriter, r *http.Request) {
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

	category, err := c.Service.Article.GetById(id)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (c *articleController) UpdateArticle(w http.ResponseWriter, r *http.Request) {
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

	var data UpdateArticleRequest
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

	if data.Category.Id == 0 {
		data.Category.CreatedBy = username
		data.Category.ModifiedBy = username
		id, err := c.Service.Category.Insert(&data.Category)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)

			return
		}

		data.Category.Id = id
	}
	article := model.Article{
		Body:       data.Body,
		Name:       data.Name,
		Url:        data.Url,
		Id:         data.Id,
		CreatedBy:  data.Category.CreatedBy,
		ModifiedBy: data.Category.ModifiedBy,
		CategoryId: data.Category.Id,
	}

	updatedArticle := c.Service.Article.Update(&article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error updating the article"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedArticle)
}

func (c *articleController) CreateNewArticle(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	w.Header().Set("Content-Type", "application/json")
	var data CreateNewArticleRequest

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

	if data.Category.Id == 0 {
		data.Category.CreatedBy = username
		data.Category.ModifiedBy = username
		id, err := c.Service.Category.Insert(&data.Category)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)

			return
		}
		data.Category.Id = id
	}

	article := model.Article{
		Body:       data.Body,
		Name:       data.Name,
		Url:        data.Url,
		CreatedBy:  username,
		ModifiedBy: username,
		CategoryId: data.Category.Id,
	}

	result, err := c.Service.Article.Insert(&article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving a new article"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
