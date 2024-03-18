CREATE PROCEDURE PR_GitHubCopilotApprovalRequests_Insert
	@GitHubCopilotId INT,
	@RequestId INT
AS
BEGIN
	INSERT INTO GitHubCopilotApprovalRequests
	(
		GitHubCopilotId,
        RequestId
	)
	VALUES(
		@GitHubCopilotId,
        @RequestId
	)
END