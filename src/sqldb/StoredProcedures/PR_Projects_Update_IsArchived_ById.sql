CREATE PROCEDURE [dbo].[PR_Projects_Update_IsArchived_ById]
  @Id Int,
	@IsArchived BIT
AS

BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	-- Insert statements for procedure here
	UPDATE 
		[dbo].[Projects]
   SET 
		[IsArchived] = @IsArchived,
		[Modified] = GETDATE()
 WHERE  
		[Id] = @Id
END