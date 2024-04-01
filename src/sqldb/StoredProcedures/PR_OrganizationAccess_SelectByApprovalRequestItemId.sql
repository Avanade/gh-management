CREATE PROCEDURE [dbo].[PR_OrganizationAccess_SelectByApprovalRequestItemId]
    @ApprovalSystemGUID UNIQUEIDENTIFIER
AS
BEGIN
    SELECT
        DISTINCT
        OA.Id,
        RO.Id 'OrganizationId',
        RO.Name 'OrganizationName',
        U.UserPrincipalName,
        U.GitHubId,
        U.GitHubUser,
        OA.Created
    FROM [dbo].[CommunityApprovals] AS CA
    LEFT JOIN [dbo].[OrganizationAccessApprovalRequests] AS OAAR ON OAAR.RequestId = CA.Id
    LEFT JOIN [dbo].[OrganizationAccess] AS OA ON OA.Id = OAAR.OrganizationAccessId
    LEFT JOIN [dbo].[RegionalOrganizations] AS RO ON RO.Id = OA.OrganizationId
    LEFT JOIN [dbo].[Users] AS U ON U.UserPrincipalName = OA.UserPrincipalName
    WHERE CA.ApprovalSystemGUID = @ApprovalSystemGUID
END