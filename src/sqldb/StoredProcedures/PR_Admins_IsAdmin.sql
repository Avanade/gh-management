CREATE PROCEDURE [dbo].[PR_Admins_IsAdmin]
	@UserPrincipalName varchar(50)
AS

IF EXISTS (
	SELECT [UserPrincipalName]
	FROM Admins
	WHERE [UserPrincipalName] = @UserPrincipalName
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