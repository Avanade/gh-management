CREATE PROCEDURE [dbo].[usp_OrganizationAccess_Select_ByApprovalSystemGUID]
    @ApprovalSystemGUID [UNIQUEIDENTIFIER]
AS
BEGIN
  SELECT DISTINCT
    [OA].[Id],
    [RO].[Id] AS [OrganizationId],
    [RO].[Name] AS [OrganizationName],
    [U].[UserPrincipalName],
    [U].[GitHubId],
    [U].[GitHubUser],
    [OA].[Created]
  FROM [dbo].[ApprovalRequest] AS [AR]
  LEFT JOIN [dbo].[OrganizationAccessApprovalRequest] AS [OAAR] ON [OAAR].[ApprovalRequestId] = [AR].[Id]
  LEFT JOIN [dbo].[OrganizationAccess] AS [OA] ON [OA].[Id] = [OAAR].[OrganizationAccessId]
  LEFT JOIN [dbo].[RegionalOrganization] AS [RO] ON [RO].[Id] = [OA].[RegionalOrganizationId]
  LEFT JOIN [dbo].[User] AS [U] ON [U].[UserPrincipalName] = [OA].[UserPrincipalName]
  WHERE [AR].[ApprovalSystemGUID] = @ApprovalSystemGUID
END
