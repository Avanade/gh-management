CREATE PROCEDURE [dbo].[PR_RepoOwners_Select_ByRepoId] 
(
	@ProjectId INT
)
AS
BEGIN
	SELECT * FROM [dbo].[RepoOwners]
    WHERE [ProjectId] = @ProjectId
END