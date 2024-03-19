CREATE PROCEDURE [dbo].[PR_GitHubCopilotApprovalRequests_SelectByCopilotId]
    @Id INT
AS
BEGIN
    SELECT 
        CA.Id,
        CA.ApproverUserPrincipalName,
        A.Name AS [ApprovalStatus],
        CA.ApprovalRemarks,
        CA.ApprovalDate,
        CA.ApprovalDescription
    FROM [dbo].[GitHubCopilotApprovalRequests] GCAR
    LEFT JOIN CommunityApprovals CA ON GCAR.RequestId = CA.Id
    LEFT JOIN ApprovalStatus A ON A.Id = CA.ApprovalStatusId
    WHERE GCAR.GitHubCopilotId = @Id
END