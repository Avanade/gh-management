CREATE PROCEDURE [dbo].[PR_Items_Insert]
	@ApplicationModuleId uniqueidentifier,
	@ApproverEmail varchar(100),
	@Subject varchar(100),
	@Body varchar(8000)
AS
	INSERT INTO Items (
		ApplicationModuleId,
		ApproverEmail,
		[Subject],
		Body
		)
	VALUES (
		@ApplicationModuleId,
		@ApproverEmail,
		@Subject,
		@Body
	)