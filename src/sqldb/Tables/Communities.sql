﻿CREATE TABLE [dbo].[Communities]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
)
