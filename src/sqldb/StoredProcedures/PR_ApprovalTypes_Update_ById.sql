/****** Object:  StoredProcedure [dbo].[PR_ApprovalTypes_Update_ById]    Script Date: 6/26/2022 4:28:13 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
ALTER PROCEDURE [dbo].[PR_ApprovalTypes_Update_ById] 
(
	@Id INT,
	@Name VARCHAR(50),
	@ApproverUserPrincipalName VARCHAR(50),
	@IsActive BIT,
	@ModifiedBy VARCHAR(50)
)
AS
BEGIN
	UPDATE [dbo].[ApprovalTypes]
	   SET [Name] = @Name
		  ,[ApproverUserPrincipalName] = @ApproverUserPrincipalName
		  ,[IsActive] = @IsActive
		  ,[Modified] = GETDATE()
		  ,[ModifiedBy] = @ModifiedBy
	 WHERE Id = @Id
	 SELECT @Id Id
END