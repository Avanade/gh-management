CREATE PROCEDURE PR_Application_IsValid
	
	@ApplicationId UNIQUEIDENTIFIER,
	@ApplicationModuleId UNIQUEIDENTIFIER

AS

IF EXISTS (
	SELECT A.[Name] [ApplicationName], AM.[Name] [ApplicationModuleName]
	FROM Applications A
	INNER JOIN ApplicationModules AM ON A.Id = AM.ApplicationId
	WHERE
	A.IsActive = 1
	AND AM.IsActive = 1
	AND A.Id = @ApplicationId
	AND AM.Id = @ApplicationModuleId)
	
	BEGIN
		SELECT '1' [IsValid]
		return 1
	END
ELSE
	BEGIN
		SELECT '0' [IsValid]
		return 0
	END