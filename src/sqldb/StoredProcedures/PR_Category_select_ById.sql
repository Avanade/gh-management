CREATE PROCEDURE [dbo].[PR_Category_select_ById]
@Id INT
AS
BEGIN
 

SELECT [Id]
      ,[Name]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Category]
  WHERE [Id] = @Id
END