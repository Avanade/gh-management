CREATE PROCEDURE [dbo].[usp_GHCPPremiumBudgetAllocation_Insert]
  @OrganizationGitHubID [INT],
  @OrganizationGitHubLogin [VARCHAR](100),
  @UserGitHubID [INT],
  @UserGitHubLogin [VARCHAR](100),
  @Amount [INT],
  @CreatedBy [VARCHAR](100)
AS
BEGIN
    DECLARE @returnID AS [INT]

    INSERT INTO [dbo].[GHCPPremiumBudgetAllocation]
    (
        [OrganizationGitHubID],
        [OrganizationGitHubLogin],
        [UserGitHubID],
        [UserGitHubLogin],
        [Amount],
        [CreatedBy],
        [Created]
    )
    VALUES
    (
        @OrganizationGitHubID,
        @OrganizationGitHubLogin,
        @UserGitHubID,
        @UserGitHubLogin,
        @Amount,
        @CreatedBy,
        GETDATE()
    )

    SET @returnID = SCOPE_IDENTITY()

    SELECT @returnID AS [Id]
END
