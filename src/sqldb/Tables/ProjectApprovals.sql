﻿CREATE TABLE [dbo].[ProjectApprovals] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [ProjectId] [INT] NOT NULL,
    [ApprovalTypeId] [INT] NOT NULL,
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
    CONSTRAINT [FK_ProjectApprovals_Projects] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Projects]([Id]),
    CONSTRAINT [FK_ProjectApprovals_ApprovalTypes] FOREIGN KEY ([ApprovalTypeId]) REFERENCES [dbo].[ApprovalTypes]([Id]),
    CONSTRAINT [FK_ProjectApprovals_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])
)
