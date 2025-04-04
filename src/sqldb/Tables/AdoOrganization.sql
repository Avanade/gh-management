CREATE TABLE [dbo].[AdoOrganization] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [Name] [VARCHAR](50) NOT NULL,
    [Purpose] [VARCHAR](100) NOT NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
)