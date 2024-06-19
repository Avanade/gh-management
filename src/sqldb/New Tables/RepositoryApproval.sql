CREATE TABLE [dbo].[RepositoryApproval] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [RepositoryId] [INT] NOT NULL,
    [RepositoryApprovalTypeId] [INT] NOT NULL,
    [ApprovalStatusId] [INT] NOT NULL,
    [ApprovalDescription] [VARCHAR](500) NULL,
    [ApprovalRemarks] [VARCHAR](255) NULL,
    [ApprovalDate] [DATETIME] NULL,
    [ApprovalSystemGUID] [UNIQUEIDENTIFIER] NULL,
    [ApprovalSystemDateSent] [DATETIME] NULL,
    [RespondedBy] [VARCHAR](100),
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_RepositoryApproval_Repository] FOREIGN KEY ([RepositoryId]) REFERENCES [dbo].[Repository]([Id]),
    CONSTRAINT [FK_RepositoryApproval_RepositoryApprovalType] FOREIGN KEY ([RepositoryApprovalTypeId]) REFERENCES [dbo].[RepositoryApprovalType]([Id]),
    CONSTRAINT [FK_RepositoryApproval_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])
)
