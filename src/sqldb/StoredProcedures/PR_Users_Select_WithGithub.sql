CREATE PROCEDURE [dbo].[PR_Users_Select_WithGithub]
 
AS

SELECT 
		[UserPrincipalName],
		[GivenName],
		[SurName],
		[JobTitle],
		[GitHubId],
		[GitHubUser],
		[Created],
		[CreatedBy],
		[Modified],
		[ModifiedBy]
  FROM 
		[dbo].[Users]
		WHERE GitHubId IS NOT NULL