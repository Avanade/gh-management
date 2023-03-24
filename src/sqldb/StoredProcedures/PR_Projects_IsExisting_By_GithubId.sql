CREATE PROCEDURE [dbo].[PR_Projects_IsExisting_By_GithubId]
	@GithubId VARCHAR(50)
AS

IF EXISTS (
	SELECT [Name]
	FROM Projects
	WHERE [GithubId] = @GithubId
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