CREATE PROCEDURE PR_CommunityApprovalRequests_Insert
	@CommunityId INT,
	@RequestId INT
AS
BEGIN
	INSERT INTO CommunityApprovalRequests
	(
		CommunityId,
        RequestId
	)
	VALUES(
		@CommunityId,
        @RequestId
	)
END