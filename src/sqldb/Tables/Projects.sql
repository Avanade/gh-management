CREATE TABLE [dbo].[Projects]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL, 
    [CoOwner] VARCHAR(100) NULL, 
    [Description] VARCHAR(1000) NULL, 
    [ConfirmAvaIP] BIT NOT NULL DEFAULT 0, 
    [ConfirmEnabledSecurity] BIT NOT NULL DEFAULT 0, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
)