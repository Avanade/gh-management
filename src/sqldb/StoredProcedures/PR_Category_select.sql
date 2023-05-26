CREATE PROCEDURE [dbo].[PR_Category_select]
AS 
BEGIN
SELECT [Id]
      ,[Name]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Category]
END