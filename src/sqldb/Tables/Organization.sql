CREATE TABLE [dbo].[Organization] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [RegionalOrganizationId] [INT] NOT NULL,
    [ClientName] [VARCHAR](100) NOT NULL,
    [ProjectName] [VARCHAR](100) NOT NULL,
    [WBS] [VARCHAR](50) NOT NULL,
    [ApprovalStatusId] [INT] NOT NULL DEFAULT 1,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_Organization_RegionalOrganization] FOREIGN KEY ([RegionalOrganizationId]) REFERENCES [dbo].[RegionalOrganization]([Id])
)
