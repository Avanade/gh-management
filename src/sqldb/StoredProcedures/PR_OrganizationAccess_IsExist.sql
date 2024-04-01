CREATE PROCEDURE [dbo].[PR_OrganizationAcess_IsExist] (
    @UserPrincipalName VARCHAR(100),
	@OrganizationId INT
) AS
BEGIN
	SELECT 
        COUNT(*) IsExist 
    FROM 
        OrganizationAccess 
    WHERE 
        UserPrincipalName = @UserPrincipalName AND OrganizationId = @OrganizationId
END