CREATE PROCEDURE [dbo].[usp_GuidanceCategory_Insert]
	@Name [VARCHAR](50),
	@CreatedBy [VARCHAR](50),
	@ModifiedBy [VARCHAR](50),
	@Id [INT] = NULL
AS
BEGIN   
    SET NOCOUNT ON 
	DECLARE @returnID AS [INT]
	SELECT @Id = [Id] FROM  [GuidanceCategory] WHERE [Name] = @Name
	  
	IF ( @Id=0  OR @Id IS NULL )
		BEGIN
			INSERT INTO [dbo].[GuidanceCategory]
			(
				[Name],
				[Created],
				[CreatedBy],
				[Modified],
				[ModifiedBy]
			)
			VALUES
			(
				@Name,
				GETDATE(),
				@CreatedBy,
				GETDATE(),
				@ModifiedBy
			)

			SET @returnID = SCOPE_IDENTITY()
			SELECT @returnID AS [Id]
		END
	ELSE
		BEGIN 
			EXEC [dbo].[usp_GuidanceCategory_Update] @Id, @Name, @CreatedBy, @ModifiedBy
			SELECT @Id AS [Id]
		END
END
