CREATE TABLE [dbo].[RepositoryApprover] (
    [RepositoryApprovalTypeId] [INT] NOT NULL,
    [ApproverUserPrincipalName] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_RepositoryApprover] PRIMARY KEY ([RepositoryApprovalTypeId], [ApproverUserPrincipalName]),
    CONSTRAINT [FK_RepositoryApprover_RepositoryApprovalType] FOREIGN KEY ([RepositoryApprovalTypeId]) REFERENCES [dbo].[RepositoryApprovalType]([Id]),
    CONSTRAINT [FK_RepositoryApprover_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
)
