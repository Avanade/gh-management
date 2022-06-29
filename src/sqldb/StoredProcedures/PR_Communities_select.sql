create PROCEDURE [dbo].[PR_Communities_select]
as 
begin

SELECT [Id]
      ,[Name]
      ,[Url]
      ,[Description]
      ,[Notes]
      ,[TradeAssocId]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Communities]



end
