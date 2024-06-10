CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Update_IsArchived]
	@Id [INT],
	@Name [VARCHAR](50),
	@ApproverUserPrincipalName [VARCHAR](50),
	@IsArchived [BIT],
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
			[ApproverUserPrincipalName] = @ApproverUserPrincipalName AND 
			[IsArchived] = 0
  	)
	BEGIN
		UPDATE
      		[dbo].[RepositoryApprovalType]
		SET
			[IsArchived] = @IsArchived,
			[Modified] = GETDATE(),
			[ModifiedBy] = @ModifiedBy
		WHERE [Id] = @Id
		
		SET @Status = 1
	END

	SELECT @Id [Id], @Status [Status]
END