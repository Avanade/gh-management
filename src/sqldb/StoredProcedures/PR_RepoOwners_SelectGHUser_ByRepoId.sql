CREATE PROCEDURE [dbo].[PR_RepoOwners_SelectGHUser_ByRepoId] 
(
	@ProjectId INT
)
AS
BEGIN
	SELECT RO.ProjectId, 
            RO.UserPrincipalName,
            U.GitHubUser
    FROM [dbo].[RepoOwners] RO
    LEFT JOIN [Users] U ON RO.UserPrincipalName = U.UserPrincipalName
    WHERE [ProjectId] = @ProjectId
    
END