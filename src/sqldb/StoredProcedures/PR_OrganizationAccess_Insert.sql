CREATE PROCEDURE [dbo].[PR_OrganizationAccess_Insert]
(
    @UserPrincipalName VARCHAR(100),
	@OrganizationId INT
) AS
BEGIN
	DECLARE @Id AS INT
 
	INSERT INTO [dbo].[OrganizationAccess]
        (
            [UserPrincipalName],
            [OrganizationId]
        )
    VALUES
        (
            @UserPrincipalName,
            @OrganizationId
        )
    SET @Id = SCOPE_IDENTITY()

    SELECT @Id AS Id
END
