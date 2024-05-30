CREATE TABLE [dbo].[GitHubCopilotApprovalRequest] (
    [GitHubCopilotId] [INT] NOT NULL,
    [ApprovalRequestId] [INT] NOT NULL,
    CONSTRAINT [PK_GitHubCopilotApprovalRequest] PRIMARY KEY ([GitHubCopilotId], [ApprovalRequestId]),
    CONSTRAINT [FK_GitHubCopilotApprovalRequest_GitHubCopilot] FOREIGN KEY ([GitHubCopilotId]) REFERENCES [dbo].[GitHubCopilot]([Id]),
    CONSTRAINT [FK_GitHubCopilotApprovalRequest_ApprovalRequest] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
