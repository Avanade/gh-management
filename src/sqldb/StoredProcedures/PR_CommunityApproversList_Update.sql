
CREATE PROCEDURE [dbo].[PR_CommunityApproversList_Update]
(
			@Id INT,
            @ApproverUserPrincipalName VARCHAR(100),
			@Disabled BIT,
            @CreatedBy VARCHAR(50),
            @ModifiedBy VARCHAR(50)
) AS
BEGIN
 
UPDATE [dbo].[CommunityApproversList]
   SET [ApproverUserPrincipalName] = @ApproverUserPrincipalName 
      ,[Created] = GETDATE()
      ,[CreatedBy] =  @CreatedBy
      ,[Modified] = GETDATE()
      ,[ModifiedBy] =@ModifiedBy
      ,[Disabled] = @Disabled
 WHERE  [Id] = @Id
END