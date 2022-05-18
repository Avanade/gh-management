ALTER PROCEDURE [dbo].[PR_Items_Insert]
	@ApplicationModuleId uniqueidentifier,
	@ApproverEmail varchar(100),
	@Subject varchar(100),
	@Body varchar(8000)
AS

	DECLARE @ResultTable table(Id [uniqueidentifier]);

	INSERT INTO Items (
		ApplicationModuleId,
		ApproverEmail,
		[Subject],
		Body
		)
	OUTPUT INSERTED.Id INTO @ResultTable
	VALUES (
		@ApplicationModuleId,
		@ApproverEmail,
		@Subject,
		@Body
	)

	SELECT [Id] FROM @ResultTable

GO