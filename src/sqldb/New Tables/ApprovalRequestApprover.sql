CREATE TABLE [dbo].[ApprovalRequestApprover] (
    [RepositoryApprovalId] [INT] NOT NULL,
    [ApproverUserPrincipalName] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_ApprovalRequestApprover] PRIMARY KEY ([RepositoryApprovalId], [ApproverUserPrincipalName]),
    CONSTRAINT [FK_ApprovalRequestApprover_RepositoryApproval] FOREIGN KEY ([RepositoryApprovalId]) REFERENCES [dbo].[RepositoryApproval]([Id]),
    CONSTRAINT [FK_ApprovalRequestApprover_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
)
