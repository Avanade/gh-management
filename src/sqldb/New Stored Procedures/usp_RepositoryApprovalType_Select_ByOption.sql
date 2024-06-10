CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Select_ByOption]
	@Offset [INT] = 0,
	@Filter [INT] = 0,
	@Search [VARCHAR](50) = '',
	@OrderBy [VARCHAR](50) = 'Name',
	@OrderType [VARCHAR](5) = 'ASC',
	@IsArchived [BIT] = 0
AS
BEGIN
	SELECT
		[Id],
		[Name],
		[ApproverUserPrincipalName],
		[IsArchived],
		[IsActive],
		[Created],
		[CreatedBy],
		[Modified],
		[ModifiedBy]
  	FROM [dbo].[RepositoryApprovalType]
	WHERE
		([Name] LIKE '%' + @Search + '%' OR
		[ApproverUserPrincipalName] LIKE '%' + @Search + '%') AND [IsArchived] = @IsArchived
	ORDER BY
		CASE WHEN @OrderType = 'ASC' THEN
			CASE @OrderBy
				WHEN 'Name' THEN [Name]
				WHEN 'ApproverUserPrincipalName' THEN [ApproverUserPrincipalName]
			END
		END,
		CASE WHEN @OrderType = 'DESC' THEN
			CASE @OrderBy
				WHEN 'Name' THEN [Name]
				WHEN 'ApproverUserPrincipalName' THEN [ApproverUserPrincipalName]
			END
		END DESC
	OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END