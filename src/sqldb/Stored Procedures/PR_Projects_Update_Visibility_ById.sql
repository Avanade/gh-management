CREATE PROCEDURE [dbo].[PR_Projects_Update_Visibility_ById]
  @Id INT,
	@VisibilityId INT
AS

BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	-- Insert statements for procedure here
	UPDATE 
		[dbo].[Projects]
   SET 
		[VisibilityId] = @VisibilityId,
		[Modified] = GETDATE()
 WHERE  
		[Id] = @Id
END