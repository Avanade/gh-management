/****** Object:  StoredProcedure [dbo].[PR_ActivityTypes_Insert]    Script Date: 6/23/2022 12:30:26 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
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
	SET @Id = (SELECT Id FROM [dbo].[ApprovalTypes] WHERE Name=@Name AND ApproverUserPrincipalName=@ApproverUserPrincipalName)

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
	END
	SELECT @Id Id
END