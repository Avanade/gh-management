CREATE PROCEDURE [dbo].[PR_ContributionAreas_Select_ByFilter](
	@Offset INT = 0,
	@Filter INT = 0,
	@Search VARCHAR(50) = '',
	@OrderBy VARCHAR(50) = 'Name',
	@OrderType VARCHAR(5) = 'ASC'
)
AS
BEGIN
	SELECT * FROM [dbo].[ContributionAreas]
	WHERE
		Name LIKE '%'+@search+'%'
	  ORDER BY
		CASE WHEN @OrderType = 'ASC' THEN Name
		END,
		CASE WHEN @OrderType = 'DESC' THEN Name
		END DESC
		OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END