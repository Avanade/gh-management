CREATE PROCEDURE [dbo].[PR_Users_IsExisting]
	@UserPrincipalName VARCHAR(100)
AS

IF EXISTS (
	SELECT UserPrincipalName
	FROM Users
	WHERE UserPrincipalName = @UserPrincipalName
)
	BEGIN
		SELECT '1' AS Result
		RETURN 1
	END
ELSE
	BEGIN
		SELECT '0' AS Result
		RETURN 0
	END