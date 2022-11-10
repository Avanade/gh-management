CREATE TABLE [dbo].[Communities]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL, 
    [Url] varchar(255) NULL,
    [Description] varchar(255) NULL,
    [Notes] varchar(255) NULL,
    [TradeAssocId] varchar(50) NULL,
    [IsExternal] BIT NOT NULL DEFAULT 0,
    [ApprovalStatusId] INT NOT NULL DEFAULT 1,
    [OnBoardingInstructions] [varchar](max) NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT FK_ApprovalStatus_Communities FOREIGN KEY (ApprovalStatusId) REFERENCES ApprovalStatus(Id)
)
