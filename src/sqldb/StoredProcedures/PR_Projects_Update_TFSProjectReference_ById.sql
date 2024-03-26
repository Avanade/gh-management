CREATE PROCEDURE [dbo].[PR_Projects_Update_TFSProjectReference_ById]
  	@Id INT,
	@TFSProjectReference VARCHAR(150),
	@Organization VARCHAR(150)
AS

BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	-- Insert statements for procedure here
	UPDATE 
		[dbo].[Projects]
   SET 
		[TFSProjectReference] = @TFSProjectReference,
		[Organization] = @Organization,
		[Modified] = GETDATE()
 WHERE  
		[Id] = @Id
END