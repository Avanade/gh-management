CREATE TABLE [dbo].[OrganizationAccessApprovalRequest] (
    [OrganizationAccessId] [INT] NOT NULL,
    [ApprovalRequestId] [INT] NOT NULL,
    CONSTRAINT [PK_OrganizationAccessApprovalRequest] PRIMARY KEY ([OrganizationAccessId], [ApprovalRequestId]),
    CONSTRAINT [FK_OrganizationAccessApprovalRequest_OrganizationAccess] FOREIGN KEY ([OrganizationAccessId]) REFERENCES [dbo].[OrganizationAccess]([Id]),
    CONSTRAINT [FK_OrganizationAccessApprovalRequest_ApprovalRequest] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
