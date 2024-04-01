CREATE PROCEDURE [dbo].[PR_OrganizationAccessApprovalRequest_Insert]
	@OrganizationAccessId INT,
	@RequestId INT
AS
BEGIN
	INSERT INTO [dbo].[OrganizationAccessApprovalRequests]
	(
		OrganizationAccessId,
        RequestId
	)
	VALUES(
		@OrganizationAccessId,
        @RequestId
	)
END