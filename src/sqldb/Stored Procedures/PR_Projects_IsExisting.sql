CREATE PROCEDURE [dbo].[PR_Projects_IsExisting]
	@Name VARCHAR(50)
AS

IF EXISTS (
	SELECT [Name]
	FROM Projects
	WHERE [Name] = @Name
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