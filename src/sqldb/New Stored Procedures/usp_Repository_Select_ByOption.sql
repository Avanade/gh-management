CREATE PROCEDURE [dbo].[usp_Repository_Select_ByOption]
	@Offset [INT] = 0,
	@Search [VARCHAR](50) = '',
  @Filter [VARCHAR](MAX) = '',
  @FilterType [TINYINT] = 0
AS
BEGIN
  SET NOCOUNT ON

  SELECT  [R].[Id],
          [R].[Name],
          [R].[AssetCode],
          [R].[ECATTID],
          [R].[Organization],
          [R].[Description],
          [R].[IsArchived],
          [R].[Created],
          [R].[RepositorySource],
          [R].[TFSProjectReference],
          [V].[Name] AS [Visibility],
          [R].[CoOwner],
          [R].[CreatedBy],
          (SELECT STRING_AGG([SRT].[Topic], ',') FROM [dbo].[RepositoryTopic] AS [SRT] WHERE [SRT].[RepositoryId] = [R].[Id]) AS [Topics],
          COUNT(*) AS [Score]
      FROM [dbo].[Repository] AS [R]
      LEFT JOIN [dbo].[Visibility] AS [V] ON [R].[VisibilityId] = [V].[Id]
      LEFT JOIN [dbo].[RepositoryTopic] AS [RT] ON [RT].[RepositoryId] = [R].[Id]
      INNER JOIN STRING_SPLIT(@Search, ' ') AS [SS] ON ([R].[Name] LIKE '%' + [SS].[value] + '%' OR [RT].[Topic] LIKE '%' + [SS].[value] + '%')
      WHERE 
          @FilterType = 0
        OR
          ([RT].[Topic] IN (SELECT VALUE FROM STRING_SPLIT(@Filter, ',')) AND @FilterType = 1)
        OR
          ([RT].[Topic] IS NULL AND @FilterType = 2)
      GROUP BY 
        [R].[Id],
        [R].[Name],
        [R].[AssetCode],
        [R].[ECATTID],
        [R].[Organization],
        [R].[Description],
        [R].[IsArchived],
        [R].[Created],
        [R].[RepositorySource],
        [R].[TFSProjectReference],
        [V].[Name],
        [R].[CoOwner],
        [R].[CreatedBy]
      ORDER BY
        [Score] DESC,
        [R].[Created] DESC
    OFFSET @Offset ROWS 
    FETCH NEXT 15 ROWS ONLY
END