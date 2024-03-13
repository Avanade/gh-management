CREATE TABLE [dbo].[CommunityApprovals]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [CommunityId] INT NULL, 
    [ApproverUserPrincipalName] VARCHAR(100) NOT NULL, 
    [ApprovalStatusId] INT NOT NULL, 
    [ApprovalDescription] VARCHAR(500) NULL,
    [ApprovalRemarks] VARCHAR(255) NULL,
    [ApprovalDate] DATETIME NULL, 
    [ApprovalSystemGUID] UNIQUEIDENTIFIER NULL,
    [ApprovalSystemDateSent] DATETIME NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT [FK_CommunityApprovals_Communities] FOREIGN KEY (CommunityId) REFERENCES Communities(Id), 
    CONSTRAINT [FK_CommunityApprovals_Users] FOREIGN KEY (ApproverUserPrincipalName) REFERENCES Users(UserPrincipalName), 
    CONSTRAINT [FK_CommunityApprovals_ApprovalStatus] FOREIGN KEY (ApprovalStatusId) REFERENCES ApprovalStatus(Id)
)
