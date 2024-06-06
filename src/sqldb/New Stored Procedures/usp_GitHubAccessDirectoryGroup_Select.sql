CREATE PROCEDURE [dbo].[usp_GitHubAccessDirectoryGroup_Select]
AS
BEGIN
  SELECT
    [ObjectId],
    [ADGroup]
  FROM [dbo].[GitHubAccessDirectoryGroup]
END
