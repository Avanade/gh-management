CREATE PROCEDURE [dbo].[usp_GitHubCopilot_Select_ByApprovalSystemGUID]
  @ApprovalSystemGUID [UNIQUEIDENTIFIER]
AS
BEGIN
  SELECT 
    [GC].[RegionalOrganizationId],
    [RO].[Name] AS [RegionName],
    [GC].[GitHubUsername],
    [GC].[GitHubId],
    [GC].[Id]
  FROM [dbo].[ApprovalRequest] AS [AR]
  LEFT JOIN [dbo].[GitHubCopilotApprovalRequest] AS [GCAR] ON [GCAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[GitHubCopilot] AS [GC] ON [GC].[Id] = [GCAR].[GitHubCopilotId]
  LEFT JOIN [dbo].[RegionalOrganization] AS [RO] ON [GC].[RegionalOrganizationId] = [RO].[Id]
  WHERE [AR].[ApprovalSystemGUID] = @ApprovalSystemGUID
END
