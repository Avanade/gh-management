CREATE TABLE [dbo].[UserAccess]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [ProjectId] INT NOT NULL, 
    [UserPrincipalName] VARCHAR(100) NOT NULL,
    [IsActive] BIT NOT NULL DEFAULT 1,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT [FK_UserAccess_Users] FOREIGN KEY (UserPrincipalName) REFERENCES Users(UserPrincipalName), 
    CONSTRAINT [FK_UserAccess_Projects] FOREIGN KEY (ProjectId) REFERENCES Projects(Id)
)
