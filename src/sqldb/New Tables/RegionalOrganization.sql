CREATE TABLE [dbo].[RegionalOrganization] (
    [Id] [INT] NOT NULL PRIMARY KEY,
    [Name] [VARCHAR](50) NOT NULL,
    [IsCleanUpMembersEnabled] [BIT] NOT NULL DEFAULT 0,
    [IsIndexRepoEnabled] [BIT] NOT NULL DEFAULT 0,
    [IsCopilotRequestEnabled] [BIT] NOT NULL DEFAULT 0,
    [IsAccessRequestEnabled] [BIT] NOT NULL DEFAULT 0,
    [IsEnabled] [BIT] NOT NULL DEFAULT 0
)
