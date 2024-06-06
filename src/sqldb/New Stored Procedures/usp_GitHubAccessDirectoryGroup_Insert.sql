CREATE PROCEDURE [dbo].[usp_GitHubAccessDirectoryGroup_Insert]
  @ObjectId [VARCHAR](100),
  @ADGroup [VARCHAR](100)
AS
BEGIN
	SET NOCOUNT ON;
	IF NOT EXISTS (SELECT [ADGroup] FROM [dbo].[GitHubAccessDirectoryGroup] WHERE [ObjectId] = @ObjectId)
    BEGIN
      INSERT INTO [dbo].[GitHubAccessDirectoryGroup]
      (
        [ObjectId],
        [ADGroup]
      )
      VALUES
      (
        @ObjectId,
        @ADGroup
      )
    END
  ELSE
    BEGIN
      UPDATE [dbo].[GitHubAccessDirectoryGroup] 
        SET [ADGroup] = @ADGroup
        WHERE [ObjectId] = @ObjectId
    END
END
