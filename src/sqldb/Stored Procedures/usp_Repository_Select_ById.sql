CREATE PROCEDURE [dbo].[usp_Repository_Select_ById]
	@Id [INT]
AS
BEGIN
	SET NOCOUNT ON;

  SELECT 
    [R].[Id],
    [GithubId],
    [R].[Name],
    [R].[Organization],
    [CoOwner],
    [Description],
    [ConfirmAvaIP],
    [ConfirmEnabledSecurity],
    [ApprovalStatusId],
    [IsArchived],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy],
    [TFSProjectReference],
    [RepositorySource],
    [V].[Name] AS [Visibility],
    (SELECT STRING_AGG([SRT].[Topic], ',') FROM [dbo].[RepositoryTopic] AS [SRT] WHERE [SRT].[RepositoryId] = [R].[Id]) AS [Topics]
  FROM [dbo].[Repository] AS [R]
  LEFT JOIN [dbo].[Visibility] AS [V] ON [R].[VisibilityId] = [V].[Id]
  WHERE [R].[Id] = @Id
END