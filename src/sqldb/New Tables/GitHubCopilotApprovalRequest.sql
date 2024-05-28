CREATE TABLE [dbo].[GitHubCopilotApprovalRequest] (
    [GitHubCopilotId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_GitHubCopilotApprovalRequest] PRIMARY KEY ([GitHubCopilotId], [RequestId]),
    CONSTRAINT [FK_GitHubCopilotApprovalRequest_GitHubCopilot] FOREIGN KEY ([GitHubCopilotId]) REFERENCES [dbo].[GitHubCopilot]([Id]),
    CONSTRAINT [FK_GitHubCopilotApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
