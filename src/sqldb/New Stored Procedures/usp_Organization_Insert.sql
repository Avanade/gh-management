CREATE PROCEDURE [dbo].[usp_Organization_Insert]
  @RegionalOrganizationId [INT],
  @ClientName [VARCHAR](100),
  @ProjectName [VARCHAR](100),
  @WBS [VARCHAR](50),
  @CreatedBy [VARCHAR](100)
AS
BEGIN
	DECLARE @returnID AS [INT]
 
	INSERT INTO [dbo].[Organization]
  (
    [RegionalOrganizationId],
    [ClientName],
    [ProjectName],
    [WBS],
    [CreatedBy],
    [Created]
  )
  VALUES
  (
    @RegionalOrganizationId,
    @ClientName,
    @ProjectName,
    @WBS,
    @CreatedBy,
    GETDATE()
  )
  SET @returnID = SCOPE_IDENTITY()

  SELECT @returnID AS [Id]
END
