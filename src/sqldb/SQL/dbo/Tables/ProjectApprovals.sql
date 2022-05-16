CREATE TABLE [dbo].[ProjectApprovals]
(
	[Id] INT NOT NULL PRIMARY KEY, 
    [ProjectId] INT NOT NULL, 
    [ApprovalTypeId] INT NOT NULL, 
    [ApproverUsername] VARCHAR(100) NOT NULL, 
    [ApprovalStatusId] INT NOT NULL, 
    [ApprovalDescription] VARCHAR(500) NULL, 
    [ApprovalDate] DATETIME NULL, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
    CONSTRAINT [FK_ProjectApprovals_Projects] FOREIGN KEY (ProjectId) REFERENCES Projects(Id), 
    CONSTRAINT [FK_ProjectApprovals_ApprovalTypes] FOREIGN KEY (ApprovalTypeId) REFERENCES ApprovalTypes(Id), 
    CONSTRAINT [FK_ProjectApprovals_ApprovalStatus] FOREIGN KEY (ApprovalStatusId) REFERENCES ApprovalStatus(Id)
)
