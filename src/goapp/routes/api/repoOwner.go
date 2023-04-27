package routes

import (
	"net/http"

	"fmt"
	ghmgmt "main/pkg/ghmgmtdb"
)

func ProjectToRepoOwner(w http.ResponseWriter, r *http.Request) {

	ProjectOwnersForRepoOwners := ghmgmt.GetProjectForRepoOwner()

	for _, ProjectOwnerForRepoOwner := range ProjectOwnersForRepoOwners {

		RepoOwners, _ := ghmgmt.RepoOwnersByUserAndProjectId(ProjectOwnerForRepoOwner.Id, ProjectOwnerForRepoOwner.UserPrincipalName)
		if len(RepoOwners) < 1 {
			err := ghmgmt.RepoOwnersInsert(ProjectOwnerForRepoOwner.Id, ProjectOwnerForRepoOwner.UserPrincipalName)
			if err != nil {

				fmt.Println(err)
			}
		}
	}

}
