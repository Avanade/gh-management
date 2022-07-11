package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
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
		fmt.Println(body)
		return
	}
	fmt.Println(body)
	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}
	db, _ := sql.Init(cp)
	switch r.Method {
	case "POST":
		param := map[string]interface{}{

			"Name":       body.Name,
			"CreatedBy":  username,
			"ModifiedBy": username,
			"Id":         body.Id,
		}

		result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Category_Insert", param)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
		fmt.Println("id")
		fmt.Println(id)
		if err != nil {
			fmt.Println(err)
		}
		for _, c := range body.CategoryArticles {
			fmt.Println(c)
			CategoryArticles := map[string]interface{}{

				"Id":          0,
				"Name ":       c.Name,
				"Url":         c.Url,
				"Body":        c.Body,
				"CategoryId ": id,
				"CreatedBy":   username,
				"ModifiedBy":  username,
			}
			fmt.Println("CategoryArticles")
			fmt.Println(CategoryArticles)
			_, err := db.ExecuteStoredProcedure("dbo.PR_CategoryArticles_Insert", CategoryArticles)
			if err != nil {
				fmt.Println(err)

			}

		}
	case "GET":
		param := map[string]interface{}{

			"Id": body.Id,
		}
		_, err := db.ExecuteStoredProcedure("dbo.PR_Communities_select_byID", param)
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

	Communities, err := db.ExecuteStoredProcedureWithResult("PR_Category_select", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(Communities)
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
	CategoryArticles, err := db.ExecuteStoredProcedureWithResult("PR_CategoryArticles_select_ById", params)
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
