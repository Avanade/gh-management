CREATE PROCEDURE [dbo].[PR_Projects_Count_OwnedPrivateRepo]
(
	@UserPrincipalName VARCHAR(100),
    @Visibility INT = 1,
    @RepositorySource VARCHAR(50) = 'GitHub',
    @Organization VARCHAR(50)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    SELECT 
        COUNT(*) AS 'OwnedPrivateRepo'
    FROM [dbo].[Projects] AS P 
        INNER JOIN [dbo].[RepoOwners] AS RO ON P.Id = RO.ProjectId 
        WHERE RO.UserPrincipalName = @UserPrincipalName 
                AND P.VisibilityId = @Visibility
                AND P.RepositorySource = @RepositorySource 
                AND P.Organization = @Organization;
END