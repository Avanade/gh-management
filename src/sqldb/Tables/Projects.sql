﻿CREATE TABLE [dbo].[Projects]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(100) NOT NULL, 
    [GithubId] INT, 
    [CoOwner] VARCHAR(100) NULL, 
    [Description] VARCHAR(MAX) NULL, 
    [ConfirmAvaIP] BIT NOT NULL DEFAULT 0, 
    [ConfirmEnabledSecurity] BIT NOT NULL DEFAULT 0, 
    [ConfirmNotClientProject] BIT NOT NULL DEFAULT 0, 
    [ApprovalStatusId] INT DEFAULT 1,
    [IsArchived] BIT NOT NULL DEFAULT 0,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL,
    [OSSsponsor] [varchar](50) NULL,
    [OSSContributionSponsorId] INT NULL,
	[Avanadeofferingsassets] [varchar](50) NULL,
	[Willbecommercialversion] [varchar](50) NULL,
	[OSSContributionInformation] [varchar](1000) NULL,
	[Newcontribution] [varchar](50) NULL,
    [VisibilityId] INT NOT NULL DEFAULT 1,
    [RepositorySource] VARCHAR(15) DEFAULT 'GitHub',
    [AssetCode] VARCHAR(50) NULL,
    [TFSProjectReference] VARCHAR(1000) NULL,
    [AssetUrl] VARCHAR(1000) NULL,
    [MaturityRating] VARCHAR(20) NULL,
    [ECATTReference] VARCHAR(1000) NULL,
    [ECATTID] INT NULL
    CONSTRAINT FK_ApprovalStatus_Projects FOREIGN KEY (ApprovalStatusId) REFERENCES ApprovalStatus(Id),
    CONSTRAINT FK_Projects_Visibility FOREIGN KEY (VisibilityId) REFERENCES Visibility(Id)
    CONSTRAINT FK_Projects_OSSContributionSponsors FOREIGN KEY (OSSContributionSponsorId) REFERENCES OSSContributionSponsors(Id)
)