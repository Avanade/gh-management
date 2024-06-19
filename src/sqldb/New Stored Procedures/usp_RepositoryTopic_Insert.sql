CREATE PROCEDURE [dbo].[usp_RepositoryTopic_Insert]
  @Topic [VARCHAR](100),
  @RepositoryId [INT]
AS
BEGIN
  INSERT INTO [dbo].[RepositoryTopic]
  (
    [Topic],
    [RepositoryId]
  )
  VALUES
  (
    @Topic,
    @RepositoryId
  )
END