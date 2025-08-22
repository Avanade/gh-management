CREATE TABLE [dbo].[GHCPPremiumBudgetAllocation] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [OrganizationGitHubID] [INT] NOT NULL,
    [OrganizationGitHubLogin] [VARCHAR](100) NOT NULL,
    [UserGitHubID] [INT] NOT NULL,
    [UserGitHubLogin] [VARCHAR](100) NOT NULL,
    [Amount] [INT] NOT NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL
)
