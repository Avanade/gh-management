CREATE PROCEDURE [dbo].[usp_Repository_Update_ECATTID]
  @Id [INT],
  @ECATTID [INT],
  @ModifiedBy [VARCHAR](100)
AS
BEGIN
	SET NOCOUNT ON;

  UPDATE 
      [dbo].[Repository]
    SET 
      [ECATTID] = @ECATTID,
      [Modified] = GETDATE(),
      [ModifiedBy] = @ModifiedBy
  WHERE  
      [Id] = @Id
END