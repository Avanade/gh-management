CREATE PROCEDURE [dbo].[PR_OrganizationAccessApprovalRequests_SelectByOrganizationAccessId]
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
    FROM [dbo].[OrganizationAccessApprovalRequests] OAAR
    LEFT JOIN CommunityApprovals CA ON OAAR.RequestId = CA.Id
    LEFT JOIN ApprovalStatus A ON A.Id = CA.ApprovalStatusId
    WHERE OAAR.OrganizationAccessId = @Id
END