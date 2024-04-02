CREATE PROCEDURE [dbo].[PR_GitHubCopilot_Insert]
(
    @Region INT,
    @GitHubId INT,
    @GitHubUsername VARCHAR(100),
    @Username VARCHAR(100)
) AS
BEGIN
	DECLARE @returnID AS INT
 
	INSERT INTO [dbo].[GitHubCopilot]
        ([Region]
        ,[GitHubId]
        ,[GitHubUsername]
        ,[CreatedBy]
        ,[Created])
    VALUES
        (@Region
        ,@GitHubId
        ,@GitHubUsername
        ,@Username
        ,GETDATE())
    SET @returnID = SCOPE_IDENTITY()

    SELECT @returnID Id
END
