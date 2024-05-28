-- This file contains SQL statements that will be executed after the build script.
/* INITIAL DATA FOR APPROVAL STATUS */
SET IDENTITY_INSERT ApprovalStatus ON IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 1
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (1, 'New') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 2
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (2, 'InReview') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 3
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (3, 'Rejected') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 4
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (4, 'NonCompliant') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 5
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (5, 'Approved') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 6
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (6, 'Retired') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 7
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (7, 'Archived')
SET IDENTITY_INSERT ApprovalStatus OFF
    /* INITIAL DATA FOR APPROVAL TYPES */
SET IDENTITY_INSERT ApprovalTypes ON IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalTypes]
        WHERE [Id] = 1
    )
INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name])
VALUES (1, 'Intellectual Property') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalTypes]
        WHERE [Id] = 2
    )
INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name])
VALUES (2, 'Legal') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalTypes]
        WHERE [Id] = 3
    )
INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name])
VALUES (3, 'Security')
SET IDENTITY_INSERT ApprovalTypes OFF
SET IDENTITY_INSERT Visibility ON IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[Visibility]
        WHERE [Id] = 1
    )
INSERT INTO [dbo].[Visibility] ([Id], [Name])
VALUES (1, 'Private') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[Visibility]
        WHERE [Id] = 2
    )
INSERT INTO [dbo].[Visibility] ([Id], [Name])
VALUES (2, 'Internal') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[Visibility]
        WHERE [Id] = 3
    )
INSERT INTO [dbo].[Visibility] ([Id], [Name])
VALUES (3, 'Public')
SET IDENTITY_INSERT Visibility OFF