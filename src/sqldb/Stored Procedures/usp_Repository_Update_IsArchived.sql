CREATE PROCEDURE [dbo].[usp_Repository_Update_IsArchived]
  @Id [INT],
	@IsArchived [BIT]
AS
BEGIN
	SET NOCOUNT ON;

  UPDATE 
		[dbo].[Repository]
  SET 
		[IsArchived] = @IsArchived,
		[Modified] = GETDATE()
  WHERE  
    [Id] = @Id
END