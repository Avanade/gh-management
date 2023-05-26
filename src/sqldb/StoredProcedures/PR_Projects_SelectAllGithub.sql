CREATE PROCEDURE [dbo].[PR_Projects_SelectAllGitHub]

AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT [Id],
       [Name],
       [TFSProjectReference]
  FROM 
       [dbo].[Projects]
WHERE
		[RepositorySource] = 'GitHub'

END