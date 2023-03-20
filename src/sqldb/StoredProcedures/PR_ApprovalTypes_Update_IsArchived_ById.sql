CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Update_IsArchived_ById] 
(
	@Id INT,
	@Name VARCHAR(50),
	@ApproverUserPrincipalName VARCHAR(50),
	@IsArchived BIT,
	@ModifiedBy BIT
)
AS
BEGIN
	DECLARE @Status AS BIT
	DECLARE @IsExist AS INT
	SET @Id = (
		SELECT Id FROM [dbo].[ApprovalTypes] WHERE 
			Id != @Id AND
			Name=@Name AND 
			ApproverUserPrincipalName=@ApproverUserPrincipalName AND 
			IsArchived = 0
	)
	SET @Status = 0

	IF @IsExist = 0
	BEGIN
		UPDATE [dbo].[ApprovalTypes]
		SET [IsArchived] = @IsArchived
			,[Modified] = GETDATE()
			,[ModifiedBy] = @ModifiedBy
		WHERE Id = @Id
		SET @Status = 1
	END

	SELECT @Id Id, @Status Status
END