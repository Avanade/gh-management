CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Select_ByApprovalStatusId]
	@ApprovalStatusId [INT]
AS
BEGIN
  SELECT
    [RA].[Id] AS [RepositoryApprovalId],
    CONVERT(varchar(36), [RA].[ApprovalSystemGUID], 1) AS [ItemId],
    [T].[Name] AS [ApprovalType],
    [R].[TFSProjectReference] AS [RepoLink],
    [R].[Name] AS [RepoName],
    [U].[Name] AS [Requester],
    [RA].[Created]
  FROM [dbo].[RepositoryApproval] AS [RA]
    INNER JOIN [dbo].[Repository] AS [R] ON [R].[Id] = [RA].[RepositoryId]
	  INNER JOIN [dbo].[RepositoryApprovalType] AS [T] ON [T].[Id] = [RA].[RepositoryApprovalTypeId]
	  INNER JOIN [dbo].[User] AS [U] ON [U].[UserPrincipalName] = [RA].[CreatedBy]
  WHERE  
	[RA].[ApprovalStatusId] = @ApprovalStatusId
END