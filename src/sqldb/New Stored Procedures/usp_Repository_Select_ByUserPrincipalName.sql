CREATE PROCEDURE [dbo].[usp_Repository_Select_ByUserPrincipalName]
	@UserPrincipalName [VARCHAR](100)
AS
BEGIN
	SET NOCOUNT ON;

  SELECT 
    [P].[Id],
    [P].[Name],
    [P].[AssetCode],
    [P].[Organization],
    [CoOwner],
    [Description],
    [ConfirmAvaIP],
    [ConfirmEnabledSecurity],
    (SELECT COUNT(*) FROM [dbo].[RepositoryApproval] WHERE [RepositoryId] = [P].[Id] AND [RespondedBy] IS NULL) AS [TotalPendingRequest],
    [ApprovalStatusId],
    [IsArchived],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy],
    [V].[Name] AS [Visibility],
    [P].[RepositorySource],
    [P].[TFSProjectReference],
    [P].[ECATTID],
    (SELECT STRING_AGG([SRT].[Topic], ',') FROM [dbo].[RepositoryTopic] AS [SRT] WHERE [SRT].[RepositoryId] = [P].[Id]) AS [Topics]
  FROM [dbo].[RepositoryOwner] AS [RO]
  LEFT JOIN [dbo].[Repository] AS [P] ON [RO].[RepositoryId] = [P].[Id]
  LEFT JOIN [dbo].[Visibility] AS [V] ON [P].[VisibilityId] = [V].[Id]
  WHERE  
   [RO].[UserPrincipalName] = @UserPrincipalName
  ORDER BY [Created] DESC
END