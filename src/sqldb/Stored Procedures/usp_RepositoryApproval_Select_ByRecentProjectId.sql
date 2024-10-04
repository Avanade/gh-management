CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Select_ByRecentProjectId]
  @Id [INT]
AS
BEGIN
  SELECT
    [RA].[Id], 
    [RA].[RepositoryId],
    [R].[Name] AS [RepositoryName],
    [RA].[RepositoryApprovalTypeId],
    [T].[Name] AS [ApprovalType],
    [RA].[ApprovalDescription],
    [S].[Name] AS [RequestStatus],
    [RA].[ApprovalDate], 
    [RA].[ApprovalRemarks]
  FROM [dbo].[RepositoryApproval] AS [RA]
    INNER JOIN [dbo].[Repository] AS [R] ON [RA].[RepositoryId] = [R].[Id]
    INNER JOIN [dbo].[RepositoryApprovalType] AS [T] ON [RA].[RepositoryApprovalTypeId] = [T].[Id]
    INNER JOIN [dbo].[ApprovalStatus] AS [S] ON [S].[Id] = [RA].[ApprovalStatusId]
  WHERE
    [RA].[Created] = (SELECT TOP(1) [Created] FROM [dbo].[RepositoryApproval] WHERE [RepositoryId] = @Id ORDER BY [Created] DESC)
  ORDER BY [RA].[Created] DESC
END
