CREATE PROCEDURE [dbo].[usp_Repository_Select_ByRepositorySource]
	@RepositorySource [VARCHAR](100) = 'GitHub'
AS
BEGIN
  SELECT 
      [Id],
      [Name],
      [GithubId],
      [TFSProjectReference]
  FROM 
      [dbo].[Repository]
  WHERE
      [RepositorySource] = @RepositorySource
END