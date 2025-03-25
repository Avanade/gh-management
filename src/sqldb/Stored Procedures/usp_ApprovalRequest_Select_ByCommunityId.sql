CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_ByCommunityId]
  @CommunityId [INT]
AS
BEGIN
  SELECT
    [AR].[Id],
    [C].[Id] AS [CommunityId],
    [C].[Name] AS [CommunityName],
    [C].[Url] AS [CommunityUrl],
    [C].[Description] AS [CommunityDescription],
    [C].[Notes] AS [CommunityNotes],
    [C].[TradeAssocId] AS [CommunityTradeAssocId],
    [C].[CommunityType] AS [CommunityType],
    [U].[Name] AS [RequesterName],
    [U].[GivenName] AS [RequesterGivenName],
    [U].[SurName] AS [RequesterSurName],
    [U].[UserPrincipalName] AS [RequesterUserPrincipalName],
    [AR].[ApproverUserPrincipalName],
    [AR].[ApprovalDescription],
    [AS].[Name] AS [ApprovalStatus]
  FROM [dbo].[CommunityApprovalRequest] AS [CAR]
    LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [CAR].[ApprovalRequestId] = [AR].[Id]
    LEFT JOIN [dbo].[Community] AS [C] ON [CAR].[CommunityId] = [C].[Id]
    LEFT JOIN [dbo].[User] AS [U] ON [C].[CreatedBy] = [U].[UserPrincipalName]
    LEFT JOIN [dbo].[ApprovalStatus] AS [AS] ON [AS].[Id] = [AR].[ApprovalStatusId]
  WHERE [CAR].[CommunityId] = @CommunityId
END