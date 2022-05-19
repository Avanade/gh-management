package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	githubAPI "main/pkg/github"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		template.UseTemplate(&w, r, "projects/new", nil)
	case "POST":

		r.ParseForm()
		fmt.Println(r)
		// func (r *Request) FormValue(key string) string

		var body models.TypNewProjectReqBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//RunExecuteStoredProcedure(body)
		//Projects_IsExisting(body)
		if Projects_IsExisting(body) {
			http.Error(w, "Existing Project Name", http.StatusBadRequest)

		} else {
			RunExecuteStoredProcedure(body)
		}

		_, err = githubAPI.CreatePrivateGitHubRepository(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
}

func RunExecuteStoredProcedure(body models.TypNewProjectReqBody) {

	cp := sql.ConnectionParam{

		ConnectionString: "Server=tcp:gh-mgmt.database.windows.net,1433;Initial Catalog=gh-mgmt;Persist Security Info=False;User ID=ghmsql;Password=Dfku2h391kj0@0cNjsl0;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;Database=GhManagementDb",
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"Name":                   body.Name,
		"CoOwner":                body.Coowner,
		"Description":            body.Description,
		"ConfirmAvaIP":           body.ConfirmAvaIP,
		"ConfirmEnabledSecurity": body.ConfirmSecIPScan,
		"CreatedBy":              body.Name,
	}
	fmt.Println(param)
	result, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Insert", param)
	if err != nil {
		fmt.Println(err)
	}

	//data, _ := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Select", nil)
	fmt.Println(result)
	//	fmt.Println(data)
}
func Projects_IsExisting(body models.TypNewProjectReqBody) bool {

	cp := sql.ConnectionParam{
		ConnectionString: "Server=tcp:gh-mgmt.database.windows.net,1433;Initial Catalog=gh-mgmt;Persist Security Info=False;User ID=ghmsql;Password=Dfku2h391kj0@0cNjsl0;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;Database=GhManagementDb",
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"Name": body.Name,
	}
	fmt.Println(param)

	// //
	// var name string
	// //rows := db.Query("select name from Projects where name = ?", body.Name).Scan(&Result)
	// row := db .QueryRow("select [Name] from Projects where [Name] = \"test\" ")
	// ////if err != nil {
	// //	fmt.Println(err)
	// //}
	// err := row.Scan(&name)
	// if err != nil {
	// 	fmt.Println(err)

	// }
	//

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_IsExisting", param)

	if err != nil {
		fmt.Println(err)
		return false
	}
	// //	var count int
	// fmt.Println("result result != nil")
	// fmt.Println(result)
	// var Result string
	// for result.Next() {

	// 	nErr := result.Scan(&Result)
	// 	if nErr != nil {
	// 		return nErr
	// 	}

	// }
	var count int
	for key, element := range result {
		//count = key
		fmt.Println("Key:", key, "=>", "Element:", element, "tets", result[element])
		count = result[element]
	}
	fmt.Println("result ", count)
	return true
	// if result == "1" {
	// 	fmt.Println("result true")
	// 	return true
	// } else {
	// 	fmt.Println("result false")
	// 	return false
	// }

	// for result.Next {
	// 	if err := result.S(&count); err != nil {
	// 		fmt.Println(err)
	// 		return false
	// 	}

	// }
	// if count > 0 {
	// 	return true
	// } else {
	// 	return false
	// }
	//data, _ := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Select", param)

}

type Result struct {
	Result string
}
