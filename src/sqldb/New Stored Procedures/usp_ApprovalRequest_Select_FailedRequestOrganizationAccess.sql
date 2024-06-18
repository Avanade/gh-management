CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_FailedRequestOrganizationAccess]
AS
BEGIN
  SELECT
    [OA].[Id] AS [Id],
    [RO].[Name] AS [RegionalOrgName],
    [UC].[GitHubUser] AS [GitHubUsername],
    [UC].[UserPrincipalName] AS [UserPrincipalName],
    STRING_AGG([AR].[ApproverUserPrincipalName], ',') AS [Approvers],
    STRING_AGG([AR].[Id], ',') AS [RequestIds]
  FROM [dbo].[ApprovalRequest] AS [AR]
  INNER JOIN [dbo].[OrganizationAccessApprovalRequest] AS [OAAR] ON [OAAR].[ApprovalRequestId] = [AR].[Id]
  INNER JOIN [dbo].[OrganizationAccess] AS [OA] ON [OA].[Id] = [OAAR].[OrganizationAccessId]
  INNER JOIN [dbo].[RegionalOrganization] AS [RO] ON [RO].[Id] = [OA].[RegionalOrganizationId]
  INNER JOIN [dbo].[User] AS [UC] ON [OA].[UserPrincipalName] = [UC].[UserPrincipalName]
  WHERE
		[AR].[ApprovalSystemGUID] IS NULL
    AND DATEDIFF(MI, [AR].[Created], GETDATE()) >=5
  GROUP BY [OA].[Id], [RO].[Name], [UC].[GitHubUser], [UC].[UserPrincipalName]
END