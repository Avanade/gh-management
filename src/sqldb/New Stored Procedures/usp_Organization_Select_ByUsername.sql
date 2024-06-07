CREATE PROCEDURE [dbo].[usp_Organization_Select_ByUsername]
  @Username [VARCHAR](100)
AS
BEGIN
  SELECT 
    [O].[Id],
    [RO].[Name],
    [O].[ClientName],
    [O].[ProjectName],
    [O].[WBS],
    [A].[Name],
    [O].[Created]
  FROM [dbo].[Organization] AS [O]
  LEFT JOIN [dbo].[RegionalOrganization] AS [RO] ON [O].[RegionalOrganizationId] = [RO].[Id]
  LEFT JOIN [dbo].[ApprovalStatus] AS [A] ON [A].[Id] = [O].[ApprovalStatusId]
  WHERE [O].[CreatedBy] = @Username
  ORDER BY [O].[Created] DESC
END
