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
  FROM [dbo].[ApprovalRequest] AR
  LEFT JOIN [dbo].[GitHubCopilotApprovalRequest] GCAR ON [GCAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[GitHubCopilot] GC ON [GC].[Id] = [GCAR].[GitHubCopilotId]
  LEFT JOIN [dbo].[RegionalOrganization] RO ON [GC].[RegionalOrganizationId] = [RO].[Id]
  WHERE [AR].[ApprovalSystemGUID] = @ApprovalSystemGUID
END
