CREATE PROCEDURE [dbo].[PR_Projects_Update_VisibilityByName]
  @Name varchar(50),
	@IsArchived BIT,
	@IsPrivate BIT,
	@ModifiedBy varchar(50)
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
		[IsPrivate] = @IsPrivate,
		[ModifiedBy] = @ModifiedBy,
		[Modified] = GETDATE()
 WHERE  
		[Name] = @Name
END