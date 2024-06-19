CREATE TABLE [dbo].[CommunityApprovalRequests] (
    [CommunityId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_CommunityApprovalRequests] PRIMARY KEY ([CommunityId], [RequestId]),
    CONSTRAINT [FK_CommunityApprovalRequests_Communities] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Communities]([Id]),
    CONSTRAINT [FK_CommunityApprovalRequests_CommunityApprovals] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[CommunityApprovals]([Id])
)
