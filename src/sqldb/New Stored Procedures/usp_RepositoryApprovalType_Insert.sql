CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Insert]
	@Name [VARCHAR](50),
	@ApproverUserPrincipalName [VARCHAR](100),
	@IsActive [BIT],
	@CreatedBy [VARCHAR](100)
AS
BEGIN
	DECLARE @Id AS [INT]
	DECLARE @Status AS [BIT]

	SET @Status = 0

	IF NOT EXISTS (
		SELECT [Id]
		FROM [dbo].[RepositoryApprovalType]
		WHERE [Name] = @Name AND [ApproverUserPrincipalName] = @ApproverUserPrincipalName AND [IsArchived] = 0
	)
	BEGIN
		INSERT INTO [dbo].[RepositoryApprovalType]
		(
			[Name],
			[ApproverUserPrincipalName],
			[IsActive],
			[IsArchived],
			[Created],
			[CreatedBy],
			[Modified],
			[ModifiedBy]
		)
		VALUES
		(
			@Name,
			@ApproverUserPrincipalName,
			@IsActive,
			0,
			GETDATE(),
			@CreatedBy,
			GETDATE(),
			@CreatedBy
		)
		SET @Id = SCOPE_IDENTITY()
		SET @Status = 1
	END
	SELECT @Id AS [Id], @Status AS [Status]
END