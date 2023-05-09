CREATE PROCEDURE [dbo].[PR_Users_Select_ByGitHubUsers]
(
	@GitHubUser VARCHAR(100)
)
 
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here

SELECT *
FROM 
	[dbo].[Users]
WHERE
    [GitHubUser] = @GitHubUser

END