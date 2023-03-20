CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Update_ById] 
(
	@Id INT,
	@Name VARCHAR(50),
	@ApproverUserPrincipalName VARCHAR(50),
	@IsActive BIT,
	@ModifiedBy VARCHAR(50)
)
AS
BEGIN
	DECLARE @Status AS BIT
	SET @Id = (SELECT Id FROM [dbo].[ApprovalTypes] WHERE Name=@Name AND ApproverUserPrincipalName=@ApproverUserPrincipalName AND IsArchived = 0)
	SET @Status = 0

	IF @Id IS NULL
	BEGIN
		UPDATE [dbo].[ApprovalTypes]
		SET [Name] = @Name
			,[ApproverUserPrincipalName] = @ApproverUserPrincipalName
			,[IsActive] = @IsActive
			,[Modified] = GETDATE()
			,[ModifiedBy] = @ModifiedBy
		WHERE Id = @Id
		
		SET @Status = 1
	END

	SELECT @Id Id, @Status Status
END