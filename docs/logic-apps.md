## Logic Apps


- **CleanupOrganization**
    - checks members of innersource organization, anyone who is not an active employee or has not used the Community Portal will be removed from the organization. 
    - checks members of the public organization, anyone who is not an active employee or has not used the Community Portal will be converted into an outside collaborator. A list of members converted into an outside collaborator is sent to OSPO. Repo admins will also receive a list of collaborators on their repositories who were converted into outside collaborators.
- **IndexOrgRepos**
    - pulls all repositories from the innersource and public organizations and updates the list of repos on database
    - pulls admins of each repository and updates repo owners table on the database
    - removes repositories that no longer exist on GitHub organizations
- **InnersourceCheckOutsideCollaborators**
    - checks outside collaborators on innersource organization and removes them
- **OpensourceCheckOutsideCollaborators**
    - checks outside collaborators on each repository and sends a list to the repo admins
- **RepoOwnserScan**
    - checks the number of admins on each repository then sends a list of repositories with less than 2 admins. Admins also receive an email encouraging them to add another admin to the repository.