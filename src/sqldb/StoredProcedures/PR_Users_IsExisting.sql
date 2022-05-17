CREATE PROCEDURE [dbo].[PR_Users_IsExisting]
	@UserPrincipalName varchar(100)
AS

IF EXISTS (
	SELECT UserPrincipalName
	FROM Users
	WHERE UserPrincipalName = @UserPrincipalName
)
	BEGIN
		SELECT '1' AS Result
		Return 1
	END
ELSE
	BEGIN
		SELECT '0' AS Result
		Return 0
	END