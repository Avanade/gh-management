CREATE PROCEDURE [dbo].[usp_AdoOrganizationApprovalRequest_Insert]
	@AdoOrganizationId [INT],
	@ApprovalRequestId [INT]
AS
BEGIN
	INSERT INTO [dbo].[AdoOrganizationApprovalRequest]
	(
		[AdoOrganizationId],
    	[ApprovalRequestId]
	)
	VALUES
	(
		@AdoOrganizationId,
    	@ApprovalRequestId
	)
END
