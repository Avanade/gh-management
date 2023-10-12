CREATE PROCEDURE [dbo].[PR_Projects_Select_ByAssetCode]
(
	@AssetCode VARCHAR(50)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT [Id],
       [GithubId],
       [Name],
       [CoOwner],
       [Description],
       [ConfirmAvaIP],
       [ConfirmEnabledSecurity],
       [ApprovalStatusId],
       [IsArchived],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy]
  FROM 
       [dbo].[Projects]
  WHERE
      [AssetCode] = @AssetCode
END