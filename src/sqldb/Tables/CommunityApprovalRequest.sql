CREATE TABLE [dbo].[CommunityApprovalRequest] (
    [CommunityId] [INT] NOT NULL,
    [ApprovalRequestId] [INT] NOT NULL,
    CONSTRAINT [PK_CommunityApprovalRequest] PRIMARY KEY ([CommunityId], [ApprovalRequestId]),
    CONSTRAINT [FK_CommunityApprovalRequest_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id]),
    CONSTRAINT [FK_CommunityApprovalRequest_ApprovalRequest] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
