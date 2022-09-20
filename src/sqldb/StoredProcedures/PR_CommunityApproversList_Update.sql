
create PROCEDURE [dbo].[PR_CommunityApproversList_Update]
(
			@Id int,
            @ApproverUserPrincipalName varchar(100),
			@Disabled bit,
            @CreatedBy varchar(50),
            @ModifiedBy varchar(50)
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
 
end