CREATE PROCEDURE [dbo].[usp_GitHubCopilot_Select_Pending]
    @Username [VARCHAR](100),
    @Organization [VARCHAR](100)
AS
BEGIN
  SELECT 
    [GC].[Id],
    [RO].[Name],
    [GC].[GitHubUsername],
    [GC].[Created]
  FROM [dbo].[GitHubCopilot] GC
  LEFT JOIN [dbo].[RegionalOrganization] AS [RO] ON [GC].[RegionalOrganizationId] = [RO].[Id]
  LEFT JOIN [dbo].[GitHubCopilotApprovalRequest] AS [GCAR] ON [GCAR].[GitHubCopilotId] = [GC].[Id]
  LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [AR].[Id] = [GCAR].[ApprovalRequestId]
  WHERE 
    [GC].[CreatedBy]=@Username AND 
    [RO].[Name] =@Organization AND 
    [AR].[ApprovalStatusId] < 3
  ORDER BY [GC].[Created] DESC
END
