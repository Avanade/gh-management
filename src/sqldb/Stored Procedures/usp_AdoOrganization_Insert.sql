CREATE PROCEDURE [dbo].[usp_AdoOrganization_Insert]
  @Name [VARCHAR](50),
  @Purpose [VARCHAR](100),
  @CreatedBy [VARCHAR](100)
AS
BEGIN
	DECLARE @returnID AS [INT]
 
	INSERT INTO [dbo].[AdoOrganization]
  (
    [Name],
    [Purpose],
    [CreatedBy],
    [Created]
  )
  VALUES
  (
    @Name,
    @Purpose,
    @CreatedBy,
    GETDATE()
  )
  SET @returnID = SCOPE_IDENTITY()

  SELECT @returnID AS [Id]
END
