CREATE TABLE [dbo].[CommunityApprovalRequest] (
    [CommunityId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_CommunityApprovalRequest] PRIMARY KEY ([CommunityId], [RequestId]),
    CONSTRAINT [FK_CommunityApprovalRequest_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id]),
    CONSTRAINT [FK_CommunityApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
