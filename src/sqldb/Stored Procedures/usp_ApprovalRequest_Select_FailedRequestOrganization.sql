CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_FailedRequestOrganization]
AS
BEGIN
  SELECT
    [O].[Id] AS [Id],
    [RO].[Id] AS [RegionId],
    [RO].[Name] AS [RegionName],
    [O].[ClientName] AS [ClientName],
    [O].[ProjectName] AS [ProjectName],
    [O].[WBS] AS [WBS],
    [O].[CreatedBy] AS [Username],
    STRING_AGG([AR].[ApproverUserPrincipalName], ',') AS [Approvers],
    STRING_AGG([AR].[Id], ',') AS [RequestIds]
  FROM [dbo].[ApprovalRequest] AS [AR]
    INNER JOIN [dbo].[OrganizationApprovalRequest] AS [OAR] ON [OAR].[ApprovalRequestId] = [AR].[Id]
    INNER JOIN [dbo].[Organization] AS [O] ON [O].[Id] = [OAR].[OrganizationId]
    INNER JOIN [dbo].[RegionalOrganization] AS [RO] ON [RO].[Id] = [O].[RegionalOrganizationId]
    INNER JOIN [dbo].[User] AS [UC] ON [O].[CreatedBy] = [UC].[UserPrincipalName]
  WHERE
		[AR].[ApprovalSystemGUID] IS NULL
    AND DATEDIFF(MI, [AR].[Created], GETDATE()) >= 5
  GROUP BY [O].[Id], [RO].[Id], [RO].[Name], [O].[ClientName], [O].[ProjectName], [O].[WBS], [O].[CreatedBy]
END