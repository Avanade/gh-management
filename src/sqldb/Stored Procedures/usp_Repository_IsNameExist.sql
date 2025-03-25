CREATE PROCEDURE [dbo].[usp_Repository_IsNameExist]
	@Name [VARCHAR](50)
AS
BEGIN
  IF EXISTS (
    SELECT [Name]
    FROM [dbo].[Repository]
    WHERE [Name] = @Name
  )
  BEGIN
		SELECT '1' AS [Result]
	END
  ELSE
	BEGIN
		SELECT '0' AS [Result]
	END
END