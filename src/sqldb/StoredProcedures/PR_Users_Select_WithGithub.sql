CREATE PROCEDURE [dbo].[PR_Users_Select_WithGithub]
 
AS

SELECT 
		[UserPrincipalName],
		[GivenName],
		[SurName],
		[JobTitle],
		[GithubUser],
		[Created],
		[CreatedBy],
		[Modified],
		[ModifiedBy]
  FROM 
		[dbo].[Users]
		WHERE GithubUser IS NOT NULL