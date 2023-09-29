CREATE TABLE [dbo].[ProjectApprovals]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [ProjectId] INT NOT NULL, 
    [ApprovalTypeId] INT NOT NULL, 
    [ApproverUserPrincipalName] VARCHAR(100), 
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
    CONSTRAINT [FK_ProjectApprovals_Projects] FOREIGN KEY (ProjectId) REFERENCES Projects(Id), 
    CONSTRAINT [FK_ProjectApprovals_ApprovalTypes] FOREIGN KEY (ApprovalTypeId) REFERENCES ApprovalTypes(Id), 
    CONSTRAINT [FK_ProjectApprovals_Users] FOREIGN KEY (ApproverUserPrincipalName) REFERENCES Users(UserPrincipalName), 
    CONSTRAINT [FK_ProjectApprovals_ApprovalStatus] FOREIGN KEY (ApprovalStatusId) REFERENCES ApprovalStatus(Id)
)
