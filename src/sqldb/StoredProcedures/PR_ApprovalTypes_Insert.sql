CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Insert]
(
	@Name VARCHAR(50),
	@ApproverUserPrincipalName VARCHAR(100),
	@IsActive BIT,
	@CreatedBy VARCHAR(100)
)
AS
BEGIN
	DECLARE @Id AS INT
	DECLARE @Status AS BIT
	SET @Id = (SELECT Id FROM [dbo].[ApprovalTypes] WHERE Name=@Name AND ApproverUserPrincipalName=@ApproverUserPrincipalName AND IsArchived = 0)
	SET @Status = 0

	IF @Id IS NULL
	BEGIN
		INSERT INTO [dbo].[ApprovalTypes] (
				Name, 
				ApproverUserPrincipalName, 
				IsActive, 
				Created, 
				CreatedBy, 
				Modified, 
				ModifiedBy
			) VALUES (
				@Name,
				@ApproverUserPrincipalName,
				@IsActive,
				getDate(),
				@CreatedBy,
				getDate(),
				@CreatedBy
			)
		SET @Id = SCOPE_IDENTITY()
		SET @Status = 1
	END
	SELECT @Id Id, @Status Status
END