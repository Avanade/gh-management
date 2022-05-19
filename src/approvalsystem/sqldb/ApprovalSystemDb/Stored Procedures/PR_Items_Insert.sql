/****** Object:  StoredProcedure [dbo].[PR_Items_Insert]    Script Date: 05/18/2022 11:24:57 pm ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
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

	SELECT dbo.UidToString(Id) [Id] FROM @ResultTable

