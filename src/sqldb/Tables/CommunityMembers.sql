CREATE TABLE [dbo].[CommunityMembers]
(
    [CommunityId] INT NOT NULL, 
    [UserPrincipalName] VARCHAR(100) NOT NULL, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT PK_CommunityMember PRIMARY KEY (CommunityId, UserPrincipalName),
    CONSTRAINT [FK_CommunityMembers_Communities] FOREIGN KEY (CommunityId) REFERENCES Communities(Id)
)
