create PROCEDURE [dbo].[PR_Category_select]
as 
begin
 

SELECT [Id]
      ,[Name]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Category]




end
