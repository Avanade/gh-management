CREATE TABLE [dbo].[OrganizationAccess]
(
    [Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [UserPrincipalName] VARCHAR(100) NOT NULL,
	[OrganizationId] INT NOT NULL,
    [Created] DATETIME NOT NULL DEFAULT GETDATE(), 
    CONSTRAINT [FK_OrganizationAccess_Users] FOREIGN KEY (UserPrincipalName) REFERENCES Users(UserPrincipalName),
    CONSTRAINT [FK_OrganizationAccess_RegionalOrganizations] FOREIGN KEY (OrganizationId) REFERENCES RegionalOrganizations(Id)
)