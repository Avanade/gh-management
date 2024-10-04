CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_FailedRequestCommunity]
AS
BEGIN
  SELECT
    [AR].[Id] AS [Id],
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
    [AR].[ApproverUserPrincipalName] AS [ApproverUserPrincipalName],
    [AR].[ApprovalDescription] AS [ApprovalDescription]
  FROM [dbo].[ApprovalRequest] AS [AR]
    INNER JOIN [dbo].[CommunityApprovalRequest] AS [CAR] ON [CAR].[ApprovalRequestId] = [AR].[Id]
    INNER JOIN [dbo].[Community] AS [C] ON [C].[Id] = [CAR].[CommunityId]
    INNER JOIN [dbo].[User] AS [U] ON [C].[CreatedBy] = [U].[UserPrincipalName]
  WHERE
		[AR].[ApprovalSystemGUID] IS NULL
    AND DATEDIFF(MI, [AR].[Created], GETDATE()) >= 5
END