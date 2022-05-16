CREATE PROCEDURE [dbo].[PR_ApplicationModules_Insert]
	@ApplicationId uniqueidentifier,
	@Name varchar(100),
	@IsActive bit = true,
	@CallbackUrl varchar(100) = NULL,
	@RequireRemarks bit = false,
	@ApprovalTypeId int = 1
AS
	INSERT INTO ApplicationModules (ApplicationId, [Name], IsActive, CallbackUrl, RequireRemarks, ApprovalTypeId)
	VALUES (@ApplicationId, @Name, @IsActive, @CallbackUrl, @RequireRemarks, @ApprovalTypeId)
