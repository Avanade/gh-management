CREATE PROCEDURE [dbo].[usp_Repository_Select_ByDateRange]
  @Start [DATETIME],
  @End [DATETIME]
AS
BEGIN
	SET NOCOUNT ON;

  SELECT 
    [Id],
    [Name],
    [Description],
    [CreatedBy]
  FROM 
    [dbo].[Repository]
  WHERE
    [Created] >= @Start AND [Created] < @End
END
