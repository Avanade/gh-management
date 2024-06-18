CREATE PROCEDURE [dbo].[PR_Communities_select]
AS
BEGIN
SELECT c.[Id]
      ,c.[Name]
      ,c.[Url]
      ,c.[Description]
      ,c.[Notes]
	  ,c.ApprovalStatusId
	  ,c.[TradeAssocId]
	  ,c.[IsExternal]
      ,c.[Created]
      ,c.[CreatedBy]
      ,c.[Modified]
      ,c.[ModifiedBy]
	  ,t.Name "ApprovalStatus"
  FROM [dbo].[Communities] c
  	INNER JOIN ApprovalStatus T ON c.ApprovalStatusId = T.Id
  WHERE 
	  c.ApprovalStatusId =5
END