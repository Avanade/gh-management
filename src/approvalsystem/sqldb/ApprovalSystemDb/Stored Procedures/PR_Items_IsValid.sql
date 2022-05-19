CREATE PROCEDURE [dbo].[PR_Items_IsValid]
	
	@ApplicationId UNIQUEIDENTIFIER,
	@ApplicationModuleId UNIQUEIDENTIFIER,
	@ItemId UNIQUEIDENTIFIER,
	@ApproverEmail varchar(100)

AS

IF EXISTS (
	SELECT I.[Subject]
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
	AND I.IsApproved IS NULL)
	
	BEGIN
		SELECT '1' [IsValid]
		return 1
	END
ELSE
	BEGIN
		SELECT '0' [IsValid]
		return 0
	END