package category

import "net/http"

type CategoryController interface {
	GetCategories(w http.ResponseWriter, r *http.Request)
	GetCategoriesById(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
}
