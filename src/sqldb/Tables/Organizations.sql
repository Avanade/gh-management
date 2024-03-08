CREATE TABLE [dbo].[Organizations]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Region] VARCHAR(100) NOT NULL,
    [ClientName] VARCHAR(100) NOT NULL,
    [ProjectName] VARCHAR(100) NOT NULL,
    [WBS] VARCHAR(50) NOT NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    CONSTRAINT [FK_Organizations_RegionalOrganizations] FOREIGN KEY (Region) REFERENCES RegionalOrganizations(Id)
)