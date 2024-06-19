CREATE PROCEDURE [dbo].[usp_RepositoryTopic_Select_PopularTopic]
  @OffSet [INT] = 0,
  @RowCount [INT] = 0
AS
BEGIN
  IF (@RowCount = 0)
  BEGIN
    SELECT
      [Topic],
      COUNT([Topic]) AS [Total]
    FROM [dbo].[RepositoryTopic]
    GROUP BY [Topic]
    ORDER BY [Total] DESC
  END
  ELSE
  BEGIN
    SELECT
      [Topic],
      COUNT([Topic]) AS [Total]
    FROM [dbo].[RepositoryTopic]
    GROUP BY [Topic]
    ORDER BY [Total] DESC
    OFFSET @OffSet ROWS
    FETCH NEXT @RowCount ROWS ONLY
  END
END