CREATE PROCEDURE [dbo].[PR_Projects_Count_OwnedRepoByVisibility]
(
	@UserPrincipalName VARCHAR(100),
    @Visibility INT = 0, -- 0 - ALL | 1 - PRIVATE | 2 - INTERNAL | 3 - PUBLIC
    @RepositorySource VARCHAR(50) = 'GitHub',
    @Organization VARCHAR(50)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    IF @Visibility = 0
        SELECT 
            COUNT(*) AS 'Total'
        FROM [dbo].[Projects] AS P 
            INNER JOIN [dbo].[RepoOwners] AS RO ON P.Id = RO.ProjectId 
            WHERE RO.UserPrincipalName = @UserPrincipalName 
                    AND P.RepositorySource = @RepositorySource 
                    AND P.Organization = @Organization;
    ELSE
        SELECT 
            COUNT(*) AS 'Total'
        FROM [dbo].[Projects] AS P 
            INNER JOIN [dbo].[RepoOwners] AS RO ON P.Id = RO.ProjectId 
            WHERE RO.UserPrincipalName = @UserPrincipalName 
                    AND P.RepositorySource = @RepositorySource 
                    AND P.Organization = @Organization
                    AND P.VisibilityId = @Visibility;
END