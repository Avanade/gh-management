CREATE TABLE [dbo].[CommunityApproversList]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [ApproverUserPrincipalName] VARCHAR(100) NOT NULL, 
    [Category] VARCHAR(100) NOT NULL DEFAULT 'community',
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL,
	[Disabled] [bit] NULL,
    CONSTRAINT [FK_CommunityApproversList_Users] FOREIGN KEY (ApproverUserPrincipalName) REFERENCES Users(UserPrincipalName)
)