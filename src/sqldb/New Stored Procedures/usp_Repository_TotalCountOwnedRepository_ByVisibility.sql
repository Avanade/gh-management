CREATE PROCEDURE [dbo].[usp_Repository_TotalCountOwnedRepository_ByVisibility]
	@UserPrincipalName [VARCHAR](100),
  @Visibility [INT] = 0, -- 0 - ALL | 1 - PRIVATE | 2 - INTERNAL | 3 - PUBLIC
  @RepositorySource [VARCHAR](50) = 'GitHub',
  @Organization [VARCHAR](50)
AS
BEGIN
	SET NOCOUNT ON;

  IF (@Visibility = 0)
  BEGIN
    SELECT 
      COUNT(*) AS [Total]
    FROM [dbo].[Repository] AS [R] 
    INNER JOIN [dbo].[RepositoryOwner] AS [RO] ON [R].[Id] = [RO].[RepositoryId] 
    WHERE [RO].[UserPrincipalName] = @UserPrincipalName 
      AND [R].[RepositorySource] = @RepositorySource 
      AND [R].[Organization] = @Organization;
  END
  ELSE
  BEGIN
    SELECT 
      COUNT(*) AS [Total]
    FROM [dbo].[Repository] AS [R] 
    INNER JOIN [dbo].[RepositoryOwner] AS [RO] ON [R].[Id] = [RO].[RepositoryId]
    WHERE [RO].[UserPrincipalName] = @UserPrincipalName 
      AND [R].[RepositorySource] = @RepositorySource 
      AND [R].[Organization] = @Organization
      AND [R].[VisibilityId] = @Visibility;
  END
END