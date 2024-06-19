CREATE TABLE [dbo].[OrganizationAccessApprovalRequests] (
    [OrganizationAccessId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_OrganizationAccessApprovalRequests] PRIMARY KEY ([OrganizationAccessId], [RequestId]),
    CONSTRAINT [FK_OrganizationAccessApprovalRequests_OrganizationAccess] FOREIGN KEY ([OrganizationAccessId]) REFERENCES [dbo].[OrganizationAccess]([Id]),
    CONSTRAINT [FK_OrganizationAccessApprovalRequests_CommunityApprovals] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[CommunityApprovals]([Id])
)
