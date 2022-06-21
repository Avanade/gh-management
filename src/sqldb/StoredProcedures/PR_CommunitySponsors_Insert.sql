Create PROCEDURE  [dbo].[PR_CommunitySponsors_Insert]
(
    -- Add the parameters for the stored procedure here
		@CommunityId int,
		@UserPrincipalName varchar(100),
		@CreatedBy varchar(50)
 
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON
DECLARE @count AS INT
	SET @count = (select MAX([CommunityId]) from [CommunitySponsors] where [UserPrincipalName] = @UserPrincipalName and[CommunityId]= @CommunityId)
	IF @count IS NULL
 BEGIN


    -- Insert statements for procedure here
INSERT INTO [dbo].[CommunitySponsors]
           ([CommunityId]
           ,[UserPrincipalName]
           ,[Created]
           ,[CreatedBy]
           ,[Modified]
           ,[ModifiedBy])
     VALUES
           (@CommunityId
           ,@UserPrincipalName
           ,GETDATE()
           ,@CreatedBy
           ,GETDATE()
           ,@CreatedBy)
 END

END
