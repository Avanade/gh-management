CREATE TABLE [dbo].[RepoOwners]
(
	[ProjectId] INT NOT NULL, 
    [UserPrincipalName] VARCHAR(100) NOT NULL, 
    CONSTRAINT PK_RepoOwner PRIMARY KEY (ProjectId, UserPrincipalName),
    CONSTRAINT FK_Projects_RepoOwners FOREIGN KEY (ProjectId) REFERENCES Projects(Id)
)