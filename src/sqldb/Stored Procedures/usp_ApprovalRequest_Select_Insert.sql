CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Select_Insert]
  @ApproverUserPrincipalName [VARCHAR](100),
  @Name [VARCHAR](485),
  @CreatedBy [VARCHAR](100)
AS
BEGIN
  DECLARE @returnID AS INT

  INSERT INTO [dbo].[ApprovalRequest]
    (
    [ApproverUserPrincipalName],
    [ApprovalStatusId],
    [ApprovalDescription],
    [CreatedBy],
    [Created],
    [ModifiedBy],
    [Modified]
    )
  VALUES(
      @ApproverUserPrincipalName,
      1,
      'For Approval - ' + @Name,
      @CreatedBy,
      GETDATE(),
      @CreatedBy,
      GETDATE()
	)

  SET @returnID = SCOPE_IDENTITY()

  SELECT @returnID AS [Id]
END