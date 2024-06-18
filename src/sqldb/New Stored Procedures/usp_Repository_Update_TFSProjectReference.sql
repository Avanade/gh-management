CREATE PROCEDURE [dbo].[usp_Repository_Update_TFSProjectReference]
	@Id [INT],
	@TFSProjectReference [VARCHAR](150),
	@Organization [VARCHAR](150)
AS
BEGIN
	SET NOCOUNT ON;

  	UPDATE 
		[dbo].[Repository]
   	SET 
		[TFSProjectReference] = @TFSProjectReference,
		[Organization] = @Organization,
		[Modified] = GETDATE()
 	WHERE  
		[Id] = @Id
END