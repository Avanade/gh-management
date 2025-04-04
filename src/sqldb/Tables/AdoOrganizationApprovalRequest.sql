CREATE TABLE [dbo].[AdoOrganizationApprovalRequest] (
    [AdoOrganizationId] [INT] NOT NULL,
    [ApprovalRequestId] [INT] NOT NULL,
    CONSTRAINT [PK_AdoOrganizationApprovalRequest] PRIMARY KEY ([AdoOrganizationId], [ApprovalRequestId]),
    CONSTRAINT [FK_AdoOrganizationApprovalRequest_AdoOrganization] FOREIGN KEY ([AdoOrganizationId]) REFERENCES [dbo].[AdoOrganization]([Id]),
    CONSTRAINT [FK_AdoOrganizationApprovalRequest_ApprovalRequest] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])
)