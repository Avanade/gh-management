CREATE TABLE [dbo].[GitHubCopilot]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Region] INT NOT NULL,
    [GitHubId] INT NOT NULL,
    [GitHubUsername] VARCHAR(100) NOT NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    CONSTRAINT [FK_GitHubCopilot_RegionalOrganizations] FOREIGN KEY (Region) REFERENCES RegionalOrganizations(Id)
)