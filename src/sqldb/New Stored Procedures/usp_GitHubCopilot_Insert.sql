CREATE PROCEDURE [dbo].[usp_GitHubCopilot_Insert]
  @RegionalOrganizationId [INT],
  @GitHubId [INT],
  @GitHubUsername [VARCHAR](100),
  @Username [VARCHAR](100)
AS
BEGIN
	DECLARE @returnID AS [INT]
 
	INSERT INTO [dbo].[GitHubCopilot]
  (
    [RegionalOrganizationId],
    [GitHubId],
    [GitHubUsername],
    [CreatedBy],
    [Created]
  )
  VALUES
  (
    @RegionalOrganizationId,
    @GitHubId,
    @GitHubUsername,
    @Username,
    GETDATE()
  )
  
  SET @returnID = SCOPE_IDENTITY()

  SELECT @returnID AS [Id]
END
