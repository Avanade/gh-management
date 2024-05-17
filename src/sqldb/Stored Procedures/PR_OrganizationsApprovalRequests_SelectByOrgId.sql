CREATE PROCEDURE [dbo].[PR_OrganizationsApprovalRequests_SelectByOrgId]
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
    FROM [dbo].[OrganizationApprovalRequests] OAR
    LEFT JOIN CommunityApprovals CA ON OAR.RequestId = CA.Id
    LEFT JOIN ApprovalStatus A ON A.Id = CA.ApprovalStatusId
    WHERE OAR.OrganizationId = @Id
END