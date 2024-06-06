CREATE PROCEDURE [dbo].[usp_GitHubCopilot_Select_ByUsername]
  @Username [VARCHAR](100)
AS
BEGIN
  SELECT 
    [GC].[Id],
    [RO].[Name],
    [GC].[GitHubUsername],
    [GC].[Created]
  FROM [dbo].[GitHubCopilot] GC
  LEFT JOIN [dbo].[RegionalOrganization] RO ON [GC].[RegionalOrganizationId] = [RO].[Id]
  WHERE [GC].[CreatedBy] = @Username
  ORDER BY [GC].[Created] DESC
END
