CREATE TABLE [dbo].[OrganizationAccessApprovalRequest] (
    [OrganizationAccessId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_OrganizationAccessApprovalRequest] PRIMARY KEY ([OrganizationAccessId], [RequestId]),
    CONSTRAINT [FK_OrganizationAccessApprovalRequest_OrganizationAccess] FOREIGN KEY ([OrganizationAccessId]) REFERENCES [dbo].[OrganizationAccess]([Id]),
    CONSTRAINT [FK_OrganizationAccessApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
