CREATE PROCEDURE [dbo].[PR_GitHubCopilot_SelectPendingByUserAndOrganization]
    @Username VARCHAR(100),
    @Organization VARCHAR(100)
AS
BEGIN
    SELECT 
        GC.Id,
        RO.Name,
        GC.GitHubUsername,
        GC.Created
    FROM [dbo].[GitHubCopilot] GC
    LEFT JOIN RegionalOrganizations RO ON GC.Region = RO.Id
    LEFT JOIN GitHubCopilotApprovalRequests GCAR ON GCAR.GitHubCopilotId = GC.Id
    LEFT JOIN CommunityApprovals CA ON CA.Id = GCAR.RequestId
    WHERE GC.CreatedBy=@Username AND RO.Name =@Organization AND CA.ApprovalStatusId < 3
    ORDER BY GC.Created DESC
END