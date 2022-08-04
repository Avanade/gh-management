CREATE TABLE [dbo].[ActivityTypes]
(
  [Id] INT NOT NULL PRIMARY KEY IDENTITY,
  [Name] VARCHAR(100) NOT NULL,
  [Created] DATETIME NOT NULL DEFAULT getdate(), 
  [CreatedBy] VARCHAR(100) NULL, 
  [Modified] DATETIME NOT NULL DEFAULT getdate(), 
  [ModifiedBy] VARCHAR(100) NULL
)
