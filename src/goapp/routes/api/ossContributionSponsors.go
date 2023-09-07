package routes

import (
	"encoding/json"
	"log"
	"net/http"

	db "main/pkg/ghmgmtdb"
)

type OssContributionSponsor struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"isArchived"`
}

func GetAllOssContributionSponsors(w http.ResponseWriter, r *http.Request) {

	sponsors, err := db.SelectAllSponsors()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(sponsors)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllEnabledOssContributionSponsors(w http.ResponseWriter, r *http.Request) {
	param := map[string]interface{}{
		"IsArchived": false,
	}

	sponsors, err := db.SelectSponsorsByIsArchived(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(sponsors)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func AddSponsor(w http.ResponseWriter, r *http.Request) {
	var data OssContributionSponsor

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"Name":       data.Name,
		"IsArchived": data.IsArchived,
	}

	_, err = db.InsertSponsor(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateSponsor(w http.ResponseWriter, r *http.Request) {
	var data OssContributionSponsor

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"Id":         data.Id,
		"Name":       data.Name,
		"IsArchived": data.IsArchived,
	}

	_, err = db.UpdateSponsor(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
