CREATE PROCEDURE [dbo].[PR_Projects_Select_AllWithMakePublicRequest]
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT *
  FROM 
       [dbo].[Projects]
  WHERE
      [ApprovalStatusId] != 1 AND [RepositorySource] = 'GitHub' AND [OSSContributionSponsorId] IS NULL
END