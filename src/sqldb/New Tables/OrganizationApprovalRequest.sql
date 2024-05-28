CREATE TABLE [dbo].[OrganizationApprovalRequest] (
    [OrganizationId] [INT] NOT NULL,
    [RequestId] [INT] NOT NULL,
    CONSTRAINT [PK_OrganizationApprovalRequest] PRIMARY KEY ([OrganizationId], [RequestId]),
    CONSTRAINT [FK_OrganizationApprovalRequest_Organization] FOREIGN KEY ([OrganizationId]) REFERENCES [dbo].[Organization]([Id]),
    CONSTRAINT [FK_OrganizationApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)
