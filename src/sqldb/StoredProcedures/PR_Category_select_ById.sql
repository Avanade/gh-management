
create PROCEDURE [dbo].[PR_Category_select_ById]
@Id int
as 
begin
 

SELECT [Id]
      ,[Name]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Category]
  where [Id] = @Id




end