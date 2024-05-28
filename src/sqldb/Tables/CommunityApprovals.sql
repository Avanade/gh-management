﻿CREATE TABLE [dbo].[CommunityApprovals] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [CommunityId] [INT] NULL,
    [ApproverUserPrincipalName] [VARCHAR](100) NOT NULL,
    [ApprovalStatusId] [INT] NOT NULL,
    [ApprovalDescription] [VARCHAR](500) NULL,
    [ApprovalRemarks] [VARCHAR](255) NULL,
    [ApprovalDate] [DATETIME] NULL,
    [ApprovalSystemGUID] [UNIQUEIDENTIFIER] NULL,
    [ApprovalSystemDateSent] [DATETIME] NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_CommunityApprovals_Users] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[Users]([UserPrincipalName]),
    CONSTRAINT [FK_CommunityApprovals_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])
)
