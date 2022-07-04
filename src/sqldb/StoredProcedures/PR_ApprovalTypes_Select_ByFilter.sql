/****** Object:  StoredProcedure [dbo].[PR_ApprovalTypes_Select_ByFilter]    Script Date: 04/07/2022 9:45:25 pm ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Select_ByFilter](
	@Offset int = 0,
	@Filter int = 0,
	@Search varchar(50) = '',
	@OrderBy varchar(50) = 'Name',
	@OrderType varchar(5) = 'ASC'
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