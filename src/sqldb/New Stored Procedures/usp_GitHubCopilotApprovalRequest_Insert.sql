CREATE PROCEDURE [dbo].[usp_GitHubCopilotApprovalRequest_Insert]
	@GitHubCopilotId [INT],
	@ApprovalRequestId [INT]
AS
BEGIN
	INSERT INTO [dbo].[GitHubCopilotApprovalRequest]
	(
		[GitHubCopilotId],
    [ApprovalRequestId]
	)
	VALUES
  (
		@GitHubCopilotId,
    	@ApprovalRequestId
	)
END
