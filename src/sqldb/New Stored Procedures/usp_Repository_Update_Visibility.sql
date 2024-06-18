CREATE PROCEDURE [dbo].[usp_Repository_Update_Visibility]
  @Id [INT],
	@VisibilityId [INT]
AS

BEGIN
	SET NOCOUNT ON;

	UPDATE 
    [dbo].[Repository]
  SET
    [VisibilityId] = @VisibilityId,
    [Modified] = GETDATE()
  WHERE  
    [Id] = @Id;
END