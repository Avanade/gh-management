CREATE PROCEDURE [dbo].[PR_Projects_IsExisting]
	@Name varchar(50)
AS

IF EXISTS (
	SELECT [Name]
	FROM Projects
	WHERE [Name] = @Name
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