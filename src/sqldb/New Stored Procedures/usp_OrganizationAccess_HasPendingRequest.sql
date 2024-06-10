CREATE PROCEDURE [dbo].[usp_OrganizationAccess_HasPendingRequest]
    @UserPrincipalName [VARCHAR](100),
	@OrganizationId [INT]
AS
BEGIN
    DECLARE @HasPendingRequest [BIT]

    SET @HasPendingRequest = (
        SELECT
            [HasPendingRequest] = CASE 
                WHEN COUNT(*) > 0 THEN 1
                ELSE 0
            END
        FROM
            [dbo].[OrganizationAccessApprovalRequest] AS [OAAR]
            LEFT JOIN [dbo].[ApprovalRequest] AS [AR] ON [OAAR].[ApprovalRequestId] = [AR].[Id]
        WHERE
            [OAAR].[OrganizationAccessId] = (
                SELECT 
                    TOP(1) Id 
                FROM 
                    [dbo].[OrganizationAccess] 
                WHERE 
                    [UserPrincipalName] = @UserPrincipalName AND 
                    [RegionalOrganizationId] = @OrganizationId
                ORDER BY Created DESC
            ) 
            AND NOT([AR].[ApprovalStatusId] NOT IN (1, 2)))

    SELECT @HasPendingRequest AS [HasPendingRequest]
END