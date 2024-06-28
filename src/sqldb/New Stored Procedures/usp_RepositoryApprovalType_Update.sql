CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Update]
	@Id [INT],
	@Name [VARCHAR](50),
	@IsActive [BIT],
	@ModifiedBy [VARCHAR](50)
AS
BEGIN
	DECLARE @Status AS [BIT]

	SET @Status = 0
	IF NOT EXISTS (
		SELECT * 
		FROM [dbo].[RepositoryApprovalType] 
		WHERE 
			[Id] != @Id AND
			[Name] = @Name AND 
			[IsArchived] = 0
	)
	BEGIN
		UPDATE
      		[dbo].[RepositoryApprovalType]
		SET 
			[Name] = @Name,
			[IsActive] = @IsActive,
			[Modified] = GETDATE(),
			[ModifiedBy] = @ModifiedBy
		WHERE [Id] = @Id
		
		SET @Status = 1
	END

	SELECT @Id AS [Id], @Status AS [Status]
END