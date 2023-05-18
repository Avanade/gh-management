
CREATE PROCEDURE [dbo].[PR_Communities_select_IManage] (@UserPrincipalName varchar(100))
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
INNER JOIN [dbo].[CommunityMembers] CM ON c.Id = CM.CommunityId
WHERE CM.UserPrincipalName = @UserPrincipalName
END