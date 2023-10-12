package routes

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"

	db "main/pkg/ghmgmtdb"

	"github.com/gorilla/mux"
)

type SearchResultItem struct {
	Code        string `xml:"Code"`
	Description string `xml:"Description"`
	Title       string `xml:"Title"`
	Type        string `sml:"Type"`
}

type ArrayOfSearchResultItem struct {
	SearchResultItem []SearchResultItem `xml:"SearchResultItem"`
	XMLNSI           string             `xml:"xmlns:i,attr"`
	XMLNS            string             `xml:"xmlns,attr"`
}

func LegacySearchHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	searchText := req["searchText"]

	param := map[string]interface{}{
		"searchText": searchText,
	}

	result, err := db.LegacySearch(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var searchResult []SearchResultItem
	for _, v := range result {
		d := SearchResultItem{
			Code:        v["Code"].(string),
			Description: v["Description"].(string),
			Title:       v["Title"].(string),
			Type:        v["Type"].(string),
		}
		searchResult = append(searchResult, d)
	}

	var finalResult ArrayOfSearchResultItem
	finalResult.SearchResultItem = searchResult
	finalResult.XMLNSI = "http://www.w3.org/2001/XMLSchema-instance"
	finalResult.XMLNS = "***REMOVED***"

	// Wraps the response to Response struct
	response, err := xml.MarshalIndent(finalResult, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	// Write
	w.Write(response)
}

func RedirectAsset(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	assetCode := req["assetCode"]

	result := db.GetProjectByAssetCode(assetCode)
	if len(result) > 0 {
		url := result[0]["TFSProjectReference"].(string)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		githubId, err := strconv.Atoi(assetCode)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = db.GetProjectByGithubId(int64(githubId))
		if len(result) > 0 {
			url := result[0]["TFSPRojectReference"].(string)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
