DELETE  
	PA 
FROM 
	[dbo].[ProjectApprovals] PA INNER JOIN [dbo].[ApprovalTypes] AT ON PA.ApprovalTypeId = AT.Id
WHERE
	AT.ApproverUserPrincipalName IS NULL;

DELETE FROM [dbo].[ApprovalTypes] WHERE ApproverUserPrincipalName IS NULL;

CREATE TABLE [dbo].[ApprovalTypes]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL,
    [ApproverUserPrincipalName] VARCHAR(100) NOT NULL,
    [IsArchived] BIT NOT NULL DEFAULT 0,
    [IsActive] BIT NOT NULL DEFAULT 1,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT FK_ApprovalTypes_Users FOREIGN KEY (ApproverUserPrincipalName) REFERENCES Users(UserPrincipalName)
)
