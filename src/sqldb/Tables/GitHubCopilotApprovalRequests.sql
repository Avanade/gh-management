CREATE TABLE [dbo].[GitHubCopilotApprovalRequests] (
    [GitHubCopilotId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_GitHubCopilotApprovalRequests] PRIMARY KEY ([GitHubCopilotId], [RequestId]),
    CONSTRAINT [FK_GitHubCopilotApprovalRequests_GitHubCopilot] FOREIGN KEY ([GitHubCopilotId]) REFERENCES [dbo].[GitHubCopilot]([Id]),
    CONSTRAINT [FK_GitHubCopilotApprovalRequests_CommunityApprovals] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[CommunityApprovals]([Id])
)
