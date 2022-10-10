CREATE PROCEDURE [dbo].[PR_Projects_Select_ByUserPrincipalName]
(
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT [p].[Id],
       [p].[Name],
       [CoOwner],
       [Description],
       [ConfirmAvaIP],
       [ConfirmEnabledSecurity],
       [ApprovalStatusId],
       [IsArchived],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy],
       [ApprovalStatusId],
       [v].[Name] AS 'Visibility',
       [p].[RepositorySource]
  FROM 
       [dbo].[Projects] AS p
  LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
  WHERE  
       [CreatedBy] = @UserPrincipalName
  ORDER BY [Created] DESC
END
