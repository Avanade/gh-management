CREATE TABLE [dbo].[RegionalOrganization] (
    [Id] [INT] NOT NULL PRIMARY KEY,
    [Name] [VARCHAR](50) NOT NULL,
    [IsRegionalOrganization] [BIT] NOT NULL DEFAULT 1,
    [IsCleanUpMembersEnabled] [BIT] NOT NULL DEFAULT 1,
    [IsIndexRepoEnabled] [BIT] NOT NULL DEFAULT 1,
    [IsCopilotRequestEnabled] [BIT] NOT NULL DEFAULT 1,
    [IsAccessRequestEnabled] [BIT] NOT NULL DEFAULT 1,
    [IsEnabled] [BIT] NOT NULL DEFAULT 1,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
	[CreatedBy] [VARCHAR](100) NULL,
	[Modified] [DATETIME] NULL,
	[ModifiedBy] [VARCHAR](100) NULL
)
