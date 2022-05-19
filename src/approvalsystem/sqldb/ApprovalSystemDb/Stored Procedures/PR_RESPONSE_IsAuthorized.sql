CREATE PROCEDURE [dbo].[PR_RESPONSE_IsAuthorized]
	
	@ApplicationId UNIQUEIDENTIFIER,
	@ApplicationModuleId UNIQUEIDENTIFIER,
	@ItemId UNIQUEIDENTIFIER,
	@ApproverEmail varchar(100)

AS

IF EXISTS (
	SELECT A.[Name] [ApplicationName], AM.[Name] [ApplicationModuleName]
	FROM Applications A
	INNER JOIN ApplicationModules AM ON A.Id = AM.ApplicationId
	INNER JOIN Items I ON AM.Id = I.ApplicationModuleId
	WHERE
	A.IsActive = 1
	AND AM.IsActive = 1
	AND A.Id = @ApplicationId
	AND AM.Id = @ApplicationModuleId
	AND I.Id = @ItemId
	AND I.ApproverEmail = @ApproverEmail
	)
	
	BEGIN
		SELECT '1' [IsAuthorized], I.IsApproved, AM.RequireRemarks
		FROM Items I
		INNER JOIN ApplicationModules AM ON I.ApplicationModuleId = AM.Id
		WHERE I.Id = @ItemId
		return 1
	END
ELSE
	BEGIN
		SELECT '0' [IsAuthorized]
		return 0
	END