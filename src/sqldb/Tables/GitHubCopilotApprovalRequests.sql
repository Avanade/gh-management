CREATE TABLE [dbo].[GitHubCopilotApprovalRequests]
(
	[GitHubCopilotId] INT NOT NULL, 
    [RequestId] INT NOT NULL,
    PRIMARY KEY (GitHubCopilotId, RequestId),
    CONSTRAINT [FK_GitHubCopilotApprovalRequests_GitHubCopilot] FOREIGN KEY (GitHubCopilotId) REFERENCES GitHubCopilot(Id),
    CONSTRAINT [FK_GitHubCopilotApprovalRequests_CommunityApprovals] FOREIGN KEY (RequestId) REFERENCES CommunityApprovals(Id)
)
