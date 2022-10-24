CREATE PROCEDURE [dbo].[PR_Communities_select_byID]
@Id INT
AS 
BEGIN
SELECT [Id]
      ,[Name]
      ,[Url]
      ,[Description]
      ,[Notes]
      ,[TradeAssocId]
      ,[IsExternal]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
      ,[ApprovalStatusId]
  FROM [dbo].[Communities]
  WHERE [Id] = @id
END
