CREATE TABLE [dbo].[Projects]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL, 
    [CoOwner] VARCHAR(100) NULL, 
    [Description] VARCHAR(1000) NULL, 
    [ConfirmAvaIP] BIT NOT NULL DEFAULT 0, 
    [ConfirmEnabledSecurity] BIT NOT NULL DEFAULT 0, 
    [ApprovalStatusId] INT DEFAULT 1,
    [IsArchived] BIT NOT NULL DEFAULT 0,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL,
    [OSSsponsor] [varchar](50) NULL,
	[Avanadeofferingsassets] [varchar](50) NULL,
	[Willbecommercialversion] [varchar](50) NULL,
	[OSSContributionInformation] [varchar](50) NULL,
	[Newcontribution] [varchar](50) NULL,
    [VisibilityId] INT NOT NULL DEFAULT 1
    CONSTRAINT FK_ApprovalStatus_Projects FOREIGN KEY (ApprovalStatusId) REFERENCES ApprovalStatus(Id),
    CONSTRAINT FK_Projects_Visibility FOREIGN KEY (VisibilityId) REFERENCES Visibility(Id)
)