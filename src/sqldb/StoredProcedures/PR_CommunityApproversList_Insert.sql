
CREATE PROCEDURE  [dbo].[PR_CommunityApproversList_Insert]
(
	@ApproverUserPrincipalName VARCHAR(100),
	@Disabled BIT = 0,
	@CreatedBy VARCHAR(50),
	@ModifiedBy VARCHAR(50),
	@Id  INT = NULL

	 
          
)
AS
BEGIN   
    SET NOCOUNT ON 
	DECLARE @returnID AS INT
	   SELECT @id= id FROM  CommunityApproversList WHERE ApproverUserPrincipalName = @ApproverUserPrincipalName

	IF ( @Id = 0  OR @Id IS NULL )
	BEGIN

INSERT INTO [dbo].[CommunityApproversList]
           ([ApproverUserPrincipalName]
           ,[Created]
           ,[CreatedBy]
           ,[Modified]
           ,[ModifiedBy]
           ,[Disabled])
     VALUES
           (@ApproverUserPrincipalName
           ,GETDATE()
           ,@CreatedBy
           ,GETDATE()
           ,@ModifiedBy
           ,@Disabled)

			SET @returnID = SCOPE_IDENTITY()
			SELECT @returnID Id
	END
	ELSE 
	BEGIN 
	 EXEC	[dbo].[PR_CommunityApproversList_Update] @Id , @ApproverUserPrincipalName,@Disabled  ,@CreatedBy ,@ModifiedBy

	SELECT @Id Id
	END
END

 