CREATE PROCEDURE [dbo].[PR_CategoryArticles_select]
AS
BEGIN
SELECT [Id]
      ,[Name]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[CategoryArticles]
END