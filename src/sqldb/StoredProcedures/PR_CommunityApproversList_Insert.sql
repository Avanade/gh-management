
create PROCEDURE  [dbo].[PR_CommunityApproversList_Insert]
(
	@ApproverUserPrincipalName varchar(100),
	@Disabled bit = 0,
	@CreatedBy varchar(50),
	@ModifiedBy varchar(50),
	@Id  int =null

	 
          
)
AS
BEGIN   
    SET NOCOUNT ON 
	DECLARE @returnID AS INT
	   select @id= id from  CommunityApproversList where ApproverUserPrincipalName = @ApproverUserPrincipalName

	IF ( @Id= 0  OR @Id is null )
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
	end
	ELSE 
	begin 
	 EXEC	[dbo].[PR_CommunityApproversList_Update] @Id , @ApproverUserPrincipalName,@Disabled  ,@CreatedBy ,@ModifiedBy

	SELECT @Id Id
	end
END

 