CREATE PROCEDURE [dbo].[PR_RegionalOrganizations_Insert]
(
			@Id INT,
            @Name VARCHAR(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
	IF NOT EXISTS (SELECT Name FROM RegionalOrganizations WHERE Id = @Id)
        BEGIN
            INSERT INTO [dbo].[RegionalOrganizations]
                (
                    [Id],
                    [Name]
                )
            VALUES
                (
                    @Id,
                    @Name
                )
        END
    ELSE
        BEGIN
            UPDATE [dbo].[RegionalOrganizations] 
            SET [Name] = @Name
            WHERE [Id] = @Id
        END
END