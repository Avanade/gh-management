CREATE TABLE [dbo].[OrganizationApprovalRequests] (
    [OrganizationId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_OrganizationApprovalRequests] PRIMARY KEY ([OrganizationId], [RequestId]),
    CONSTRAINT [FK_OrganizationApprovalRequests_Organizations] FOREIGN KEY ([OrganizationId]) REFERENCES [dbo].[Organizations]([Id]),
    CONSTRAINT [FK_OrganizationApprovalRequests_CommunityApprovals] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[CommunityApprovals]([Id])
)
