CREATE PROCEDURE [dbo].[PR_Users_GetEmailByGitHubId]
(
	@GithubId VARCHAR(50)
)
AS
BEGIN
SELECT 
		[UserPrincipalName]
  FROM 
		[dbo].[Users]
  WHERE  
	    [GitHubId] = @GithubId
END