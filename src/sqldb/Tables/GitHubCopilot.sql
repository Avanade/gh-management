CREATE TABLE [dbo].[GitHubCopilot] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [RegionalOrganizationId] [INT] NOT NULL,
    [GitHubId] [INT] NOT NULL,
    [GitHubUsername] [VARCHAR](100) NOT NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_GitHubCopilot_RegionalOrganization] FOREIGN KEY ([RegionalOrganizationId]) REFERENCES [dbo].[RegionalOrganization]([Id])
)
