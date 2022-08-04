CREATE TABLE [dbo].[Menu]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL, 
    [Url] VARCHAR(50) NOT NULL DEFAULT '#', 
    [IconPath] VARCHAR(50) NULL,
    [IsActive] BIT NOT NULL DEFAULT 1, 
    [ExternalLink] BIT NOT NULL DEFAULT 0, 
    [Order] INT NOT NULL DEFAULT 1, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    
)