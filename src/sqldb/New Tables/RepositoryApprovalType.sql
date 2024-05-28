﻿CREATE TABLE [dbo].[RepositoryApprovalType] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [Name] [VARCHAR](50) NOT NULL,
    [ApproverUserPrincipalName] [VARCHAR](100) NULL,
    [IsArchived] [BIT] NOT NULL DEFAULT 0,
    [IsActive] [BIT] NOT NULL DEFAULT 1,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_RepositoryApprovalType_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
)