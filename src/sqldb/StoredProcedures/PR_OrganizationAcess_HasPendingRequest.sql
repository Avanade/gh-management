CREATE PROCEDURE [dbo].[PR_OrganizationAcess_HasPendingRequest] (
    @UserPrincipalName VARCHAR(100),
	@OrganizationId INT
) AS
BEGIN
    DECLARE @HasPendingRequest BIT

    SET @HasPendingRequest = (SELECT
            HasPendingRequest = CASE 
                WHEN COUNT(*) > 0 THEN 1
                ELSE 0
            END
        FROM
            OrganizationAccessApprovalRequests AS OAAR
            LEFT JOIN CommunityApprovals AS CA ON OAAR.RequestId = CA.Id
        WHERE
            OAAR.OrganizationAccessId = (
                SELECT 
                    TOP(1) Id 
                FROM 
                    OrganizationAccess 
                WHERE 
                    UserPrincipalName = @UserPrincipalName AND 
                    OrganizationId = @OrganizationId
                ORDER BY Created DESC
            ) 
            AND NOT(CA.ApprovalStatusId NOT IN (1, 2)))

    SELECT @HasPendingRequest HasPendingRequest
END