﻿﻿CREATE TABLE [dbo].[Projects] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [Name] [VARCHAR](100) NOT NULL,
    [GithubId] [INT],
    [CoOwner] [VARCHAR](100) NULL,
    [Description] [VARCHAR](MAX) NULL,
    [ConfirmAvaIP] [BIT] NOT NULL DEFAULT 0,
    [ConfirmEnabledSecurity] [BIT] NOT NULL DEFAULT 0,
    [ConfirmNotClientProject] [BIT] NOT NULL DEFAULT 0,
    [ApprovalStatusId] [INT] DEFAULT 1,
    [IsArchived] [BIT] NOT NULL DEFAULT 0,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    [Organization] [VARCHAR](100),
    [OSSContributionSponsorId] [INT] NULL,
    [Avanadeofferingsassets] [VARCHAR](50) NULL,
    [Willbecommercialversion] [VARCHAR](50) NULL,
    [OSSContributionInformation] [VARCHAR](1000) NULL,
    [Newcontribution] [VARCHAR](50) NULL,
    [VisibilityId] [INT] NOT NULL DEFAULT 1,
    [RepositorySource] [VARCHAR](15) DEFAULT 'GitHub',
    [AssetCode] [VARCHAR](50) NULL,
    [TFSProjectReference] [VARCHAR](1000) NULL,
    [AssetUrl] [VARCHAR](1000) NULL,
    [MaturityRating] [VARCHAR](20) NULL,
    [ECATTReference] [VARCHAR](1000) NULL,
    [ECATTID] [INT] NULL,
    CONSTRAINT [FK_Projects_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id]),
    CONSTRAINT [FK_Projects_Visibility] FOREIGN KEY ([VisibilityId]) REFERENCES [dbo].[Visibility]([Id]),
    CONSTRAINT [FK_Projects_OSSContributionSponsors] FOREIGN KEY ([OSSContributionSponsorId]) REFERENCES [dbo].[OSSContributionSponsors]([Id])
)
