CREATE PROCEDURE [dbo].[usp_Repository_Select_ByGitHubId]
	@GithubId [INT]
AS
BEGIN
	SET NOCOUNT ON;

  SELECT 
    [Id],
    [GithubId],
    [Name],
    [Organization],
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
    [GithubId] = @GithubId
END