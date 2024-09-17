CREATE TABLE [dbo].[OrganizationApprovalRequest] (
    [OrganizationId] [INT] NOT NULL,
    [ApprovalRequestId] [INT] NOT NULL,
    CONSTRAINT [PK_OrganizationApprovalRequest] PRIMARY KEY ([OrganizationId], [ApprovalRequestId]),
    CONSTRAINT [FK_OrganizationApprovalRequest_Organization] FOREIGN KEY ([OrganizationId]) REFERENCES [dbo].[Organization]([Id]),
    CONSTRAINT [FK_OrganizationApprovalRequest_ApprovalRequest] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
