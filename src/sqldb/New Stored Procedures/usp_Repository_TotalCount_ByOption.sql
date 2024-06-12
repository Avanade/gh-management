CREATE PROCEDURE [dbo].[usp_Repository_TotalCount_ByOption]
  @Search [VARCHAR](50) = '',
  @Filter [VARCHAR](MAX) = '',
  @FilterType [TINYINT] = 0
AS
BEGIN
  SET NOCOUNT ON

  SELECT COUNT(*) AS [Total]
  FROM (
		SELECT COUNT(*) AS [Total]
    FROM [dbo].[Repository] AS [R]
    LEFT JOIN [dbo].[RepositoryTopic] AS [RT] ON [RT].[RepositoryId] = [R].[Id]
    INNER JOIN STRING_SPLIT(@Search, ' ') AS [S] ON (
			[R].[Name] LIKE '%' + [S].[value] + '%' OR [RT].[Topic] LIKE '%' + [S].[value] + '%'
		)
    WHERE @FilterType = 0 OR (
      [RT].[Topic] IN (
        SELECT VALUE
        FROM STRING_SPLIT(@Filter, ',')
      ) AND @FilterType = 1
    ) OR (
      [RT].[Topic] IS NULL AND @FilterType = 2
    )
    GROUP BY [Id]
	) AS [Total]
END