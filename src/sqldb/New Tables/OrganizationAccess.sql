CREATE TABLE [dbo].[OrganizationAccess] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [UserPrincipalName] [VARCHAR](100) NOT NULL,
    [RegionalOrganizationId] [INT] NOT NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    CONSTRAINT [FK_OrganizationAccess_User] FOREIGN KEY ([UserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName]),
    CONSTRAINT [FK_OrganizationAccess_RegionalOrganization] FOREIGN KEY ([RegionalOrganizationId]) REFERENCES [dbo].[RegionalOrganization]([Id])
)
