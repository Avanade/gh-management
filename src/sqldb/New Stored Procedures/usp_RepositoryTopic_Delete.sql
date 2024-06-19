CREATE PROCEDURE [dbo].[usp_RepositoryTopic_Delete]
  @RepositoryId [INT]
AS
BEGIN
  DELETE FROM [dbo].[RepositoryTopic] WHERE [RepositoryId] = @RepositoryId
END