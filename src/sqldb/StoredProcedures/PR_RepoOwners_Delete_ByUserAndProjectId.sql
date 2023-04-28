CREATE PROCEDURE [dbo].[PR_RepoOwners_Delete_ByUserAndProjectId] 
(
    @ProjectId INT,
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
	DELETE FROM [dbo].[RepoOwners]
    WHERE [ProjectId] = @ProjectId AND [UserPrincipalName] = @UserPrincipalName
END