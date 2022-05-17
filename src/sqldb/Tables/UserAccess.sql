CREATE TABLE [dbo].[UserAccess]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [ProjectId] INT NOT NULL, 
    [Username] VARCHAR(100) NOT NULL, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
    CONSTRAINT [FK_UserAccess_Users] FOREIGN KEY (Username) REFERENCES Users(Username), 
    CONSTRAINT [FK_UserAccess_Projects] FOREIGN KEY (ProjectId) REFERENCES Projects(Id)
)
