CREATE PROCEDURE [dbo].[PR_GitHubAccess_Insert]
(
			@ObjectId VARCHAR(100),
            @ADGroup VARCHAR(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
	IF NOT EXISTS (SELECT ADGroup FROM GitHubAccess WHERE ObjectId = @ObjectId)
        BEGIN
            INSERT INTO [dbo].[GitHubAccess]
                (
                    [ObjectId],
                    [ADGroup]
                )
            VALUES
                (
                    @ObjectId,
                    @ADGroup
                )
        END
    ELSE
        BEGIN
            UPDATE [dbo].[GitHubAccess] 
            SET [ADGroup] = @ADGroup
            WHERE [ObjectId] = @ObjectId
        END
END