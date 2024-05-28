CREATE TABLE [dbo].[RepositoryApprover] (
    [ApprovalTypeId] [INT] NOT NULL,
    [ApproverEmail] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_RepositoryApprover] PRIMARY KEY ([ApprovalTypeId], [ApproverEmail]),
    CONSTRAINT [FK_RepositoryApprover_RepositoryApprovalType] FOREIGN KEY ([ApprovalTypeId]) REFERENCES [dbo].[RepositoryApprovalType]([Id]),
    CONSTRAINT [FK_RepositoryApprover_User] FOREIGN KEY ([ApproverEmail]) REFERENCES [dbo].[User]([UserPrincipalName])
)
