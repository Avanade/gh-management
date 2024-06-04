CREATE PROCEDURE [dbo].[usp_ActivityType_Insert]
	@Name [VARCHAR](100)
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @Id AS INT
	SET @Id = (SELECT [Id]
	FROM [dbo].[ActivityType]
	WHERE [Name] = @Name)

	IF @Id IS NULL
	BEGIN
		INSERT INTO [dbo].[ActivityType]
			([Name])
		VALUES
			(@Name)
		SET @Id = SCOPE_IDENTITY()
	END
	SELECT @Id Id
END