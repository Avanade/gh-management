CREATE PROCEDURE [dbo].[PR_RepoOwners_Select_ByUserPrincipalName] 
(
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
	SELECT * FROM [dbo].[RepoOwners]
    WHERE [UserPrincipalName] = @UserPrincipalName
END