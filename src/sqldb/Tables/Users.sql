CREATE TABLE [dbo].[Users]
(
    [UserPrincipalName] VARCHAR(100) NOT NULL PRIMARY KEY, 
    [Name] VARCHAR(100) NOT NULL, 
    [GivenName] VARCHAR(100) NULL, 
    [SurName] VARCHAR(100) NULL, 
    [JobTitle] VARCHAR(100) NULL, 
    [GitHubId] VARCHAR(100) NULL,
    [GitHubUser] VARCHAR(100) NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
)
