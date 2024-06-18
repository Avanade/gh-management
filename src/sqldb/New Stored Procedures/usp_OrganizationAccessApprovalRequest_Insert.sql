CREATE PROCEDURE [dbo].[usp_OrganizationAccessApprovalRequest_Insert]
	@OrganizationAccessId [INT],
	@ApprovalRequestId [INT]
AS
BEGIN
	INSERT INTO [dbo].[OrganizationAccessApprovalRequest]
	(
		[OrganizationAccessId],
    [ApprovalRequestId]
	)
	VALUES(
		@OrganizationAccessId,
    @ApprovalRequestId
	)
END