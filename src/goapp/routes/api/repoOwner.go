package routes

import (
	"log"
	"net/http"

	ghmgmt "main/pkg/ghmgmtdb"
)

func InitProjectToRepoOwner(w http.ResponseWriter, r *http.Request) {

	ProjectOwnersForRepoOwners := ghmgmt.GetProjectForRepoOwner()

	for _, ProjectOwnerForRepoOwner := range ProjectOwnersForRepoOwners {

		RepoOwners, _ := ghmgmt.RepoOwnersByUserAndProjectId(ProjectOwnerForRepoOwner.Id, ProjectOwnerForRepoOwner.UserPrincipalName)
		if len(RepoOwners) < 1 {
			err := ghmgmt.RepoOwnersInsert(ProjectOwnerForRepoOwner.Id, ProjectOwnerForRepoOwner.UserPrincipalName)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

}
