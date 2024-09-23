CREATE PROCEDURE [dbo].[usp_ContributionArea_Select_ByOption]
  @Offset [INT] = 0,
  @Filter [INT] = 0,
  @Search [VARCHAR](50) = '',
  @OrderBy [VARCHAR](50) = 'Name',
  @OrderType [VARCHAR](5) = 'ASC'
AS
BEGIN
  SELECT
    [Id],
    [Name],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[ContributionArea]
  WHERE
		Name LIKE '%' + @Search + '%'
  ORDER BY
		CASE WHEN @OrderType = 'ASC' THEN Name
		END,
		CASE WHEN @OrderType = 'DESC' THEN Name
		END DESC
		OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END