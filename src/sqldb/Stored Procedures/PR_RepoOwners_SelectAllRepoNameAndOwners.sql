create PROCEDURE [dbo].[PR_RepoOwners_SelectAllRepoNameAndOwners] 
 
AS
BEGIN
select 
	RO.ProjectId ,
	P.Name, 
	Ro.UserPrincipalName 
from 
	RepoOwners RO left  join 
	Projects P on RO.ProjectId = p.Id
END

