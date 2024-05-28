CREATE TABLE [dbo].[ApprovalRequestApprover] (
    [ApprovalRequestId] [INT] NOT NULL,
    [ApproverEmail] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_ApprovalRequestApprover] PRIMARY KEY ([ApprovalRequestId], [ApproverEmail]),
    CONSTRAINT [FK_ApprovalRequestApprover_RepositoryApproval] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[RepositoryApproval]([Id]),
    CONSTRAINT [FK_ApprovalRequestApprover_User] FOREIGN KEY ([ApproverEmail]) REFERENCES [dbo].[User]([UserPrincipalName])
)