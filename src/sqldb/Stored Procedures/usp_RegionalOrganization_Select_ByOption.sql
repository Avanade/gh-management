CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select_ByOption]
  @Offset [INT] = 0,
	@Filter [INT] = 10,
	@Search [VARCHAR](50) = '',
	@OrderBy [VARCHAR](50) = 'Date',
	@OrderType [VARCHAR](5) = 'ASC'
AS
BEGIN
    SELECT
      [Id],
      [Name],
      [IsCleanUpMembersEnabled],
      [IsIndexRepoEnabled],
      [IsCopilotRequestEnabled],
      [IsAccessRequestEnabled],
      [IsEnabled],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy],
      COUNT(*) AS [Score]
    FROM 
      [dbo].[RegionalOrganization] AS [RO]
    INNER JOIN 
      STRING_SPLIT(@Search, ' ') AS [SS] ON ([RO].[Name] LIKE '%'+[SS].[value]+'%')
    GROUP BY
      [Id],
      [Name],
      [IsCleanUpMembersEnabled],
      [IsIndexRepoEnabled],
      [IsCopilotRequestEnabled],
      [IsAccessRequestEnabled],
      [IsEnabled],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy]
    ORDER BY
      [Score] DESC,
      CASE WHEN @OrderType = 'ASC' THEN [Name] END,
      CASE WHEN @OrderType = 'DESC' THEN [Name] END DESC
    OFFSET @Offset ROWS
    FETCH NEXT @Filter ROWS ONLY
END