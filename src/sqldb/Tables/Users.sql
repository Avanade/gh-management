CREATE TABLE [dbo].[Users]
(
    [UserPrincipalName] VARCHAR(100) NOT NULL PRIMARY KEY, 
    [GivenName] VARCHAR(100) NOT NULL, 
    [SurName] VARCHAR(100) NOT NULL, 
    [JobTitle] VARCHAR(100) NOT NULL, 
    [GithubUser] VARCHAR(100) NOT NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
)
