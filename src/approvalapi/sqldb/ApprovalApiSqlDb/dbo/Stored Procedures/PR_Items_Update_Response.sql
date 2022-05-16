CREATE PROCEDURE [dbo].[PR_Items_Update_Response]
	@Id uniqueidentifier,
	@IsApproved bit,
	@ApproverRemarks varchar(255)
AS
	UPDATE Items
	SET
		IsApproved = @IsApproved,
		ApproverRemarks = @ApproverRemarks,
		DateResponded = GETDATE(),
		Modified = GETDATE()