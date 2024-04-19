CREATE PROCEDURE PR_OrganizationsApprovalRequests_Insert
	@OrganizationId INT,
	@RequestId INT
AS
BEGIN
	INSERT INTO OrganizationApprovalRequests
	(
		OrganizationId,
        RequestId
	)
	VALUES(
		@OrganizationId,
        @RequestId
	)
END