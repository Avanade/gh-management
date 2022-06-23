/****** Object:  StoredProcedure [dbo].[PR_ActivityTypes_Insert]    Script Date: 6/23/2022 12:30:26 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
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
	UPDATE [dbo].[ApprovalTypes]
	   SET [Name] = @Name
		  ,[ApproverUserPrincipalName] = @ApproverUserPrincipalName
		  ,[IsActive] = @IsActive
		  ,[Modified] = GETDATE()
		  ,[ModifiedBy] = @ModifiedBy
	 WHERE Id = @Id
END