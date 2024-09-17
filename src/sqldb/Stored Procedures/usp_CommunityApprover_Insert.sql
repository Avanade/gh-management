CREATE PROCEDURE [dbo].[usp_CommunityApprover_Insert]
  @ApproverUserPrincipalName [VARCHAR](100),
  @IsDisabled [BIT] = 0,
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50),
  @Id [INT] = 0,
  @GuidanceCategory [VARCHAR](100)
AS
BEGIN   
  SET NOCOUNT ON 
  SELECT @Id = [Id] FROM  [dbo].[CommunityApprover] WHERE [ApproverUserPrincipalName] = @ApproverUserPrincipalName AND [GuidanceCategory] = @GuidanceCategory

  IF NOT EXISTS (SELECT * FROM [CommunityApprover] WHERE [Id] = @Id)
  BEGIN
    INSERT INTO [dbo].[CommunityApprover]
    (
      [ApproverUserPrincipalName],
      [GuidanceCategory],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy],
      [IsDisabled]
    )
    VALUES
    (
      @ApproverUserPrincipalName,
      @GuidanceCategory,
      GETDATE(),
      @CreatedBy,
      GETDATE(),
      @ModifiedBy,
      @IsDisabled
    )
    SET @Id = SCOPE_IDENTITY()
  END
  ELSE 
  BEGIN 
    EXEC	[dbo].[usp_CommunityApprover_Update] @Id, @ApproverUserPrincipalName, @IsDisabled , @CreatedBy, @ModifiedBy
  END

  SELECT @Id AS [Id]
END