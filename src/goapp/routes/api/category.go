package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	session "main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func CategoryAPIHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypCategory
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	// cp := sql.ConnectionParam{

	// 	ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	// }

	// db, _ := sql.Init(cp)
	switch r.Method {
	case "POST":
		param := map[string]interface{}{

			"Name":       body.Name,
			"CreatedBy":  username,
			"ModifiedBy": username,
			"Id":         body.Id,
		}

		result, err := ghmgmt.CategoryInsert(param)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))

		if err != nil {
			fmt.Println(err)
		}
		for _, c := range body.CategoryArticles {

			CategoryArticles := map[string]interface{}{

				"Id":          0,
				"Name ":       c.Name,
				"Url":         c.Url,
				"Body":        c.Body,
				"CategoryId ": id,
				"CreatedBy":   username,
				"ModifiedBy":  username,
			}

			_, err := ghmgmt.CategoryArticlesInsert(CategoryArticles)
			if err != nil {
				fmt.Println(err)

			}

		}
	case "GET":
		param := map[string]interface{}{

			"Id": body.Id,
		}
		_, err := ghmgmt.CommunitiesSelectByID(param)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func CategoryListAPIHandler(w http.ResponseWriter, r *http.Request) {
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get project list

	Communities, err := ghmgmt.CategorySelect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(Communities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)

}

func GetCategoryArticlesById(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	CategoryArticles, err := ghmgmt.CategoryArticlesselectById(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(CategoryArticles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCategoryArticlesByArticlesID(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id

	CategoryArticles, err := ghmgmt.CategoryArticlesSelectByArticlesID(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(CategoryArticles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	Category, err := ghmgmt.CategorySelectById(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func CategoryArticlesUpdate(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypCategoryArticles
	//var body models.TypCategory
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}
	db, _ := sql.Init(cp)

	param1 := map[string]interface{}{

		"Name":       body.CategoryName,
		"CreatedBy":  username,
		"ModifiedBy": username,
		"Id":         body.CategoryId,
	}

	result, err := ghmgmt.CategoryInsert(param1)
	if err != nil {
		fmt.Println(err)

	}
	id2, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	param := map[string]interface{}{
		"Id":         body.Id,
		"Name":       body.Name,
		"Url":        body.Url,
		"Body":       body.Body,
		"CategoryId": id2,
		"CreatedBy":  username,
		"ModifiedBy": username,
	}
	_, err2 := db.ExecuteStoredProcedure("dbo.PR_CategoryArticles_Update", param)
	if err2 != nil {
		fmt.Println(err)
		return
	}

}

func CategoryUpdate(w http.ResponseWriter, r *http.Request) {

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	//var body models.TypCategoryArticles
	var body models.TypCategory

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	// cp := sql.ConnectionParam{

	// 	ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	// }
	// db, err := sql.Init(cp)

	param1 := map[string]interface{}{

		"Name":       body.Name,
		"CreatedBy":  username,
		"ModifiedBy": username,
		"Id":         body.Id,
	}

	_, err2 := ghmgmt.CategoryInsert(param1)
	if err != nil {
		fmt.Println(err2)

		return
	}

}
