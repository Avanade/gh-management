CREATE PROCEDURE [dbo].[usp_OrganizationApprovalRequest_Insert]
	@OrganizationId [INT],
	@ApprovalRequestId [INT]
AS
BEGIN
	INSERT INTO [dbo].[OrganizationApprovalRequest]
	(
		[OrganizationId],
    [ApprovalRequestId]
	)
	VALUES
  (
		@OrganizationId,
    @ApprovalRequestId
	)
END
