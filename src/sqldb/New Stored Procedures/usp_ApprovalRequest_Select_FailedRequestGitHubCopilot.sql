CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_FailedRequestGitHubCopilot]
AS
BEGIN
  SELECT
    [GC].[Id] AS [Id],
    [RO].[Id] AS [RegionId],
    [RO].[Name] AS [RegionName],
    [U].[GitHubId] AS [GitHubId],
    [U].[GitHubUser] AS [GitHubUsername],
    [U].[UserPrincipalName] AS [Username],
    STRING_AGG([AR].[ApproverUserPrincipalName], ',') AS [Approvers],
    STRING_AGG([AR].[Id], ',')  AS [RequestIds]
  FROM [dbo].[ApprovalRequest] AS [AR]
    INNER JOIN [dbo].[GitHubCopilotApprovalRequest] AS [GCAR] ON [GCAR].[ApprovalRequestId] = [AR].[Id]
    INNER JOIN [dbo].[GitHubCopilot] AS [GC] ON [GC].[Id] = [GCAR].[GitHubCopilotId]
    INNER JOIN [dbo].[RegionalOrganization] AS [RO] ON [RO].[Id] = [GC].[RegionalOrganizationId]
    INNER JOIN [dbo].[User] AS [U] ON [GC].[CreatedBy] = [U].[UserPrincipalName]
  WHERE
		[AR].[ApprovalSystemGUID] IS NULL
    AND DATEDIFF(MI, [AR].[Created], GETDATE()) >= 5
  GROUP BY [GC].[Id], [RO].[Id], [RO].[Name], [U].[GitHubId], [U].[GitHubUser], [U].[UserPrincipalName]
END