CREATE PROCEDURE [dbo].[PR_Projects_Select_By_DateRange]
(
  @Start DATETIME,
  @End DATETIME
)

AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
SELECT [Id],
       [Name],
       [Description],
       [CreatedBy]
  FROM 
       [dbo].[Projects]
  WHERE
      [Created] >= @Start AND [Created] < @End

END
