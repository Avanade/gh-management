CREATE PROCEDURE [dbo].[usp_CommunityApprover_Update]
  @Id [INT],
  @ApproverUserPrincipalName [VARCHAR](100),
  @IsDisabled [BIT],
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50)
AS
BEGIN
  UPDATE [dbo].[CommunityApprover]
   SET 
    [ApproverUserPrincipalName] = @ApproverUserPrincipalName ,
    [Created] = GETDATE(),
    [CreatedBy] =  @CreatedBy,
    [Modified] = GETDATE(),
    [ModifiedBy] =@ModifiedBy,
    [IsDisabled] = @IsDisabled
  WHERE  [Id] = @Id
END