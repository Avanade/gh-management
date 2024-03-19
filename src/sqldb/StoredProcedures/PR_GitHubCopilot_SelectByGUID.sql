CREATE PROCEDURE [dbo].[PR_GitHubCopilot_SelectByGUID]
    @ApprovalSystemGUID UNIQUEIDENTIFIER
AS
BEGIN
    SELECT 
        GC.Region,
        RO.Name AS RegionName,
        GC.GitHubUsername,
        GC.GitHubId,
        GC.Id
    FROM [dbo].[CommunityApprovals] CA
    LEFT JOIN GitHubCopilotApprovalRequests GCAR ON GCAR.RequestId = CA.Id
    LEFT JOIN GitHubCopilot GC ON GC.Id = GCAR.GitHubCopilotId
    LEFT JOIN RegionalOrganizations RO ON GC.Region = RO.Id
    WHERE CA.ApprovalSystemGUID = @ApprovalSystemGUID
END