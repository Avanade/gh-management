CREATE TABLE [dbo].[Users]
(
    [Username] VARCHAR(100) NOT NULL, 
    [FirstName] VARCHAR(50) NOT NULL, 
    [LastName] VARCHAR(50) NOT NULL, 
    [Email] VARCHAR(50) NOT NULL, 
    [GithubUser] VARCHAR(100) NOT NULL, 
    PRIMARY KEY ([Username]) , 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
)
