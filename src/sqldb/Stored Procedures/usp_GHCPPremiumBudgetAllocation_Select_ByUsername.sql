CREATE PROCEDURE [dbo].[usp_GHCPPremiumBudgetAllocation_Select_ByUsername]
  @Username [VARCHAR](100)
AS
BEGIN
  SELECT
    [Id],
    [OrganizationGitHubID],
    [OrganizationGitHubLogin],
    [UserGitHubID],
    [UserGitHubLogin],
    [Amount],
    [Created],
    [CreatedBy]
  FROM [dbo].[GHCPPremiumBudgetAllocation] AS [GC]
  WHERE [GC].[CreatedBy] = @Username
END
