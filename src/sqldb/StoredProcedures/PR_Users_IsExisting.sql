CREATE PROCEDURE [dbo].[PR_Users_IsExisting]
	@Username varchar(100)
AS

IF EXISTS (
	SELECT Username
	FROM Users
	WHERE Username = @Username
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