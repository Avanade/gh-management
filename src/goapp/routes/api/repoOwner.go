package routes

import (
	"log"
	"net/http"

	db "main/pkg/ghmgmtdb"
)

func InitProjectToRepoOwner(w http.ResponseWriter, r *http.Request) {

	ProjectOwnersForRepoOwners := db.GetProjectForRepoOwner()

	for _, ProjectOwnerForRepoOwner := range ProjectOwnersForRepoOwners {

		RepoOwners, _ := db.RepoOwnersByUserAndProjectId(ProjectOwnerForRepoOwner.Id, ProjectOwnerForRepoOwner.UserPrincipalName)
		if len(RepoOwners) < 1 {
			err := db.RepoOwnersInsert(ProjectOwnerForRepoOwner.Id, ProjectOwnerForRepoOwner.UserPrincipalName)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

}
