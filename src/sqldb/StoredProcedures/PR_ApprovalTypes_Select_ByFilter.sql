CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Select_ByFilter](
	@Offset INT = 0,
	@Filter INT = 0,
	@Search VARCHAR(50) = '',
	@OrderBy VARCHAR(50) = 'Name',
	@OrderType VARCHAR(5) = 'ASC'
)
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalTypes]
	WHERE
		Name LIKE '%'+@search+'%' OR
		ApproverUserPrincipalName LIKE '%'+@search+'%'
	  ORDER BY
		CASE WHEN @OrderType = 'ASC' THEN
			CASE @OrderBy
				WHEN 'Name' THEN Name
				WHEN 'ApproverUserPrincipalName' THEN ApproverUserPrincipalName
			END
		END,
		CASE WHEN @OrderType = 'DESC' THEN
			CASE @OrderBy
				WHEN 'Name' THEN Name
				WHEN 'ApproverUserPrincipalName' THEN ApproverUserPrincipalName
			END
		END DESC
		OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END