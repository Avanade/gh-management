CREATE PROCEDURE [dbo].[PR_RepoOwners_Select_ByUserAndProjectId] 
(
    @ProjectId INT,
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
	SELECT * FROM [dbo].[RepoOwners]
    WHERE [ProjectId] = @ProjectId AND [UserPrincipalName] = @UserPrincipalName
END