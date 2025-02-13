package category

import "net/http"

type CategoryController interface {
	GetCategories(w http.ResponseWriter, r *http.Request)
	GetCategoryById(w http.ResponseWriter, r *http.Request)
	CreateCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
}
