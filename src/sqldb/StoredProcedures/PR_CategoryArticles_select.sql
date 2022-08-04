create PROCEDURE [dbo].[PR_CategoryArticles_select]
as 
begin
 

SELECT [Id]
      ,[Name]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[CategoryArticles]




end