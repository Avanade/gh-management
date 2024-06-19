CREATE PROCEDURE [dbo].[usp_CommunityApproverRequest_Insert]
	@CommunityId INT,
	@RequestId INT
AS
BEGIN
	INSERT INTO [dbo].[CommunityApprovalRequest]
	(
		[CommunityId],
    [ApprovalRequestId]
	)
	VALUES(
		@CommunityId,
    @RequestId
	)
END