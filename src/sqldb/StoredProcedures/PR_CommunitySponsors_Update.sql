
CREATE PROCEDURE [dbo].[PR_CommunitySponsors_Update]
(
    -- Add the parameters for the stored procedure here
		@CommunityId INT,
		@UserPrincipalName VARCHAR(100),
		@CreatedBy VARCHAR(50) 
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
   
UPDATE [dbo].[CommunitySponsors]
   SET [CommunityId] =@CommunityId
      ,[UserPrincipalName] = @UserPrincipalName
      ,[Created] = GETDATE()
      ,[CreatedBy] = @CreatedBy
      ,[Modified] = GETDATE()
      ,[ModifiedBy] =  @CreatedBy
 WHERE [CommunityId] =@CommunityId
END
