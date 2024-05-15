CREATE PROCEDURE [dbo].[PR_Projects_Select_ByGithubId]
(
	@GithubId INT
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
       [dbo].[Projects]
  WHERE
      [GithubId] = @GithubId
END
