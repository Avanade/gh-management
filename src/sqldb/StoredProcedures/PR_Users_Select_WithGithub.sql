CREATE PROCEDURE [dbo].[PR_Users_Select_WithGithub]
 
AS

SELECT 
		[UserPrincipalName],
		[GivenName],
		[SurName],
		[JobTitle],
		[GithubId],
		[GithubUser],
		[Created],
		[CreatedBy],
		[Modified],
		[ModifiedBy]
  FROM 
		[dbo].[Users]
		WHERE GithubId IS NOT NULL