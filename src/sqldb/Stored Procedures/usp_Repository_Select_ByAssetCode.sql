CREATE PROCEDURE [dbo].[usp_Repository_Select_ByAssetCode]
	@AssetCode [VARCHAR](50)
AS
BEGIN
	SET NOCOUNT ON;

  SELECT 
    [Id],
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
    [ModifiedBy],
    [TFSProjectReference]
  FROM 
    [dbo].[Repository]
  WHERE
    [AssetCode] = @AssetCode AND
    [VisibilityId] != 1
END