CREATE PROCEDURE [dbo].[usp_Community_Select_MyAllApproved]
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  SELECT 
    [C].[Id],
    [C].[Name],
    [C].[Url],
    [C].[Description],
    [C].[Notes],
    [C].[ApprovalStatusId],
    [C].[TradeAssocId],
    [C].[IsExternal],
    [C].[Created],
    [C].[CreatedBy],
    [C].[Modified],
    [C].[ModifiedBy],
    [T].[Name] AS [ApprovalStatus]
  FROM [dbo].[Community] AS [C]
  INNER JOIN [dbo].[ApprovalStatus] AS [T] ON [C].[ApprovalStatusId] = [T].[Id]
  INNER JOIN [dbo].[CommunityMember] AS [CM] ON [C].Id = [CM].[CommunityId]
  WHERE [CM].[UserPrincipalName] = @UserPrincipalName AND [C].[ApprovalStatusId] = 5
END