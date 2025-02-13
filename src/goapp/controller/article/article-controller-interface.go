package article

import "net/http"

type ArticleController interface {
	GetArticlesByCategoryId(w http.ResponseWriter, r *http.Request)
	GetArticleById(w http.ResponseWriter, r *http.Request)
	UpdateArticle(w http.ResponseWriter, r *http.Request)
	CreateNewArticle(w http.ResponseWriter, r *http.Request)
}
