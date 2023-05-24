CREATE PROCEDURE [dbo].[PR_Projects_Select_ByRepositorySource]
(
	@RepositorySource VARCHAR(100) = 'GitHub'
)
AS
BEGIN
    SELECT 
        [Id],
        [Name],
        [GithubId]
    FROM 
        [dbo].[Projects]
    WHERE
        RepositorySource=@RepositorySource
END