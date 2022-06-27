/****** Object:  StoredProcedure [dbo].[PR_ApprovalTypes_Select]    Script Date: 6/24/2022 1:49:08 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
ALTER PROCEDURE [dbo].[PR_ApprovalTypes_Select_ByFilter](
	@Offset int = 0,
	@Filter int = 0,
	@Search varchar(50) = '',
	@OrderBy varchar(5) = 'ASC'
)
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalTypes]
	WHERE
		Name LIKE '%'+@search+'%' OR
		ApproverUserPrincipalName LIKE '%'+@search+'%'
	  ORDER BY
		CASE WHEN @OrderBy='ASC' THEN Modified  END,
		CASE WHEN @OrderBy='DESC' THEN Modified  END DESC
		OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END