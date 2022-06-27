create PROCEDURE [dbo].[PR_Communities_select_byID]
@Id int
as 
begin

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
  FROM [dbo].[Communities]
  where [Id] = @Id


end
